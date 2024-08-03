package usecase

import (
	"context"
	"errors"
	dtos "github.com/dafailyasa/invoice-item-service/internal/dtos/invoice"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/internal/repository"
	elasticrepository "github.com/dafailyasa/invoice-item-service/internal/repository/es"
	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type InvoiceUseCase struct {
	DB                 *gorm.DB
	CustomerRepository *repository.CustomerRepository
	ItemRepository     *repository.ItemRepository
	InvoiceRepository  *repository.InvoiceRepository
	ElasticRepository  *elasticrepository.ElasticRepository
	InvoiceItemUseCase *InvoiceItemUseCase
	ItemUseCase        *ItemUseCase
}

func NewInvoiceUseCase(
	db *gorm.DB,
	customerRepository *repository.CustomerRepository,
	itemRepository *repository.ItemRepository,
	invoiceRepository *repository.InvoiceRepository,
	elasticRepository *elasticrepository.ElasticRepository,
	invoiceItemUseCase *InvoiceItemUseCase,
	itemUseCase *ItemUseCase,
) *InvoiceUseCase {
	return &InvoiceUseCase{
		DB:                 db,
		CustomerRepository: customerRepository,
		InvoiceRepository:  invoiceRepository,
		ItemRepository:     itemRepository,
		ElasticRepository:  elasticRepository,
		InvoiceItemUseCase: invoiceItemUseCase,
		ItemUseCase:        itemUseCase,
	}
}

func (uc *InvoiceUseCase) Create(ctx context.Context, request *dtos.CreateOrUpdateInvoiceRequest) (*entities.Invoice, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	c := new(entities.Customer)
	if err := uc.CustomerRepository.FindById(ctx, tx, c, request.CustomerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Invoice{}, apperror.BadRequest(apperror.ErrCustomerNotFound)
		}
		return &entities.Invoice{}, apperror.InternalServerError(err)
	}
	if err := uc.ItemUseCase.ValidateItemsID(ctx, tx, request.Items); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(apperror.ErrItemNotFound)
	}

	invArgs := request.MapToEntityInvoice()
	if err := uc.InvoiceRepository.Create(tx, &invArgs); err != nil {
		return &entities.Invoice{}, apperror.InternalServerError(err)
	}

	var (
		invItemsArgs []entities.InvoiceItem
		totalAmount  float64
	)

	for _, item := range request.Items {
		d := item.MapToEntityInvoiceItem(invArgs.ID)
		totalAmount = totalAmount + d.Amount
		invItemsArgs = append(invItemsArgs, d)
	}

	invArgs.TotalAmount = totalAmount
	invArgs.ItemCount = len(invItemsArgs)
	invArgs.InvoiceItems = invItemsArgs

	if err := uc.InvoiceRepository.Update(tx, &invArgs); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	res, err := uc.InvoiceRepository.FindInvoiceWithRelations(ctx, tx, invArgs.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Invoice{}, apperror.BadRequest(apperror.ErrInvNotFound)
		}
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	if err := uc.ProcessToES(ctx, &res); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	if err := tx.Commit().Error; err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	return dtos.MapToInvoiceDetailResponse(&res), nil
}
func (uc *InvoiceUseCase) Update(ctx context.Context, request *dtos.CreateOrUpdateInvoiceRequest, invID string) (*entities.Invoice, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	i := new(entities.Invoice)
	if err := uc.InvoiceRepository.FindById(ctx, tx, i, invID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Invoice{}, apperror.BadRequest(apperror.ErrInvNotFound)
		}
	}

	c := new(entities.Customer)
	if err := uc.CustomerRepository.FindById(ctx, tx, c, request.CustomerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Invoice{}, apperror.BadRequest(apperror.ErrCustomerNotFound)
		}
		return &entities.Invoice{}, apperror.InternalServerError(err)
	}

	if err := uc.ItemUseCase.ValidateItemsID(ctx, tx, request.Items); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(apperror.ErrItemNotFound)
	}

	invArgs := request.MapToEntityInvoiceUpdate(i)

	if err := uc.InvoiceRepository.Update(tx, &invArgs); err != nil {
		return &entities.Invoice{}, apperror.InternalServerError(err)
	}

	var (
		invItemsArgs []entities.InvoiceItem
		totalAmount  float64
	)

	for _, item := range request.Items {
		d := item.MapToEntityInvoiceItem(invArgs.ID)
		totalAmount = totalAmount + d.Amount
		invItemsArgs = append(invItemsArgs, d)
	}

	invItemIds := lo.Map(invItemsArgs, func(item entities.InvoiceItem, _ int) uint64 {
		return item.ItemID
	})

	if err := uc.InvoiceItemUseCase.RemoveByIds(ctx, tx, i.ID, invItemIds); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	invArgs.TotalAmount = totalAmount
	invArgs.ItemCount = len(invItemsArgs)
	invArgs.InvoiceItems = invItemsArgs

	if err := uc.InvoiceRepository.Update(tx, &invArgs); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	res, err := uc.InvoiceRepository.FindInvoiceWithRelations(ctx, tx, invArgs.ID)
	if err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	if err := uc.ProcessToES(ctx, &res); err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	if err := tx.Commit().Error; err != nil {
		return &entities.Invoice{}, apperror.BadRequest(err)
	}

	return &res, nil
}
func (uc *InvoiceUseCase) ProcessToES(ctx context.Context, res *entities.Invoice) error {
	invArgsES := entities.MapToElasticInvoiceEntity(res)
	if err := uc.ElasticRepository.Index(ctx, *invArgsES); err != nil {
		return apperror.BadRequest(err)
	}

	return nil
}

func (uc *InvoiceUseCase) Search(ctx context.Context, pagination *pagination.PaginationRequest) ([]entities.ElasticInvoice, error) {
	i, err := uc.ElasticRepository.Search(ctx, pagination)
	if err != nil {
		return nil, apperror.BadRequest(err)
	}

	return i, nil
}
