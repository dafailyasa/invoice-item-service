package usecase

import (
	"context"
	"errors"
	dtos2 "github.com/dafailyasa/invoice-item-service/internal/dtos/invoice-item"
	dtos "github.com/dafailyasa/invoice-item-service/internal/dtos/item"
	"github.com/sourcegraph/conc/pool"
	"sync"

	"github.com/dafailyasa/invoice-item-service/internal/dtos/item/converter"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/internal/repository"
	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	"gorm.io/gorm"
)

type ItemUseCase struct {
	DB             *gorm.DB
	ItemRepository *repository.ItemRepository
}

func NewItemUseCase(db *gorm.DB, itemRepository *repository.ItemRepository) *ItemUseCase {
	return &ItemUseCase{
		DB:             db,
		ItemRepository: itemRepository,
	}
}

func (uc *ItemUseCase) Create(ctx context.Context, request *dtos.CreateOrUpdateItemRequest) (*entities.Item, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	args := request.MapToEntityItem()

	if err := uc.ItemRepository.Create(tx, &args); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &entities.Item{}, apperror.BadRequest(apperror.ErrItemNameExist)
		}
		return &entities.Item{}, apperror.BadRequest(err)
	}

	if err := tx.Commit().Error; err != nil {
		return &entities.Item{}, apperror.BadRequest(err)
	}

	return converter.ItemToResponse(&args), nil

}
func (uc *ItemUseCase) Update(ctx context.Context, request *dtos.CreateOrUpdateItemRequest, id string) (*entities.Item, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	i := new(entities.Item)
	if err := uc.ItemRepository.FindById(ctx, tx, i, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Item{}, apperror.BadRequest(apperror.ErrItemNotFound)
		}
		return &entities.Item{}, apperror.InternalServerError(err)
	}

	i.Description = request.Description
	i.Name = request.Name

	if err := uc.ItemRepository.Update(tx, i); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &entities.Item{}, apperror.BadRequest(apperror.ErrItemNameExist)
		}
		return &entities.Item{}, apperror.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return &entities.Item{}, apperror.BadRequest(err)
	}

	return converter.ItemToResponse(i), nil
}
func (uc *ItemUseCase) List(ctx context.Context, pagination *pagination.PaginationRequest) (customers []entities.Item, err error) {
	i, err := uc.ItemRepository.Search(ctx, uc.DB, pagination)
	if err != nil {
		return nil, apperror.BadRequest(err)
	}
	return i, nil
}
func (uc *ItemUseCase) ValidateItemsID(ctx context.Context, tx *gorm.DB, items []dtos2.InvoiceItem) error {
	p := pool.New().WithContext(ctx).WithMaxGoroutines(10)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var collectedErrors []error

	for _, item := range items {
		item := item
		wg.Add(1)

		p.Go(func(ctx context.Context) error {
			defer wg.Done()

			i := new(entities.Item)
			if err := uc.ItemRepository.FindById(ctx, tx, i, item.ItemID); err != nil {
				var customErr error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					customErr = apperror.BadRequest(apperror.ErrItemNotFound)
				} else {
					customErr = apperror.InternalServerError(apperror.ErrCustomerNotFound)
				}

				mu.Lock()
				collectedErrors = append(collectedErrors, customErr)
				mu.Unlock()

				return customErr
			}
			return nil
		})
	}

	if err := p.Wait(); err != nil {
		return err
	}

	if len(collectedErrors) > 0 {
		return collectedErrors[0]
	}

	return nil
}
