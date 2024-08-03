package usecase

import (
	"context"
	"errors"
	dtos "github.com/dafailyasa/invoice-item-service/internal/dtos/customer"
	"github.com/dafailyasa/invoice-item-service/internal/dtos/customer/converter"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/internal/repository"

	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	"gorm.io/gorm"
)

type CustomerUseCase struct {
	DB                 *gorm.DB
	CustomerRepository *repository.CustomerRepository
}

func NewCustomerUseCase(db *gorm.DB, customerRepository *repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		DB:                 db,
		CustomerRepository: customerRepository,
	}
}

func (uc *CustomerUseCase) Create(ctx context.Context, request *dtos.CreateOrUpdateCustomerRequest) (*entities.Customer, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	args := request.MapToEntityCustomer()
	if err := uc.CustomerRepository.Create(tx, &args); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &entities.Customer{}, apperror.BadRequest(apperror.ErrCustomerEmailExist)
		}
		return &entities.Customer{}, apperror.BadRequest(err)
	}

	if err := tx.Commit().Error; err != nil {
		return &entities.Customer{}, apperror.BadRequest(err)
	}

	return converter.CustomerToResponse(&args), nil
}
func (uc *CustomerUseCase) Detail(ctx context.Context, id string) (*entities.Customer, error) {
	c := new(entities.Customer)
	if err := uc.CustomerRepository.FindById(ctx, uc.DB, c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Customer{}, apperror.BadRequest(apperror.ErrCustomerNotFound)
		}
		return &entities.Customer{}, apperror.InternalServerError(err)
	}

	return converter.CustomerToResponse(c), nil
}
func (uc *CustomerUseCase) List(ctx context.Context, pagination *pagination.PaginationRequest) (customers []entities.Customer, err error) {
	c, err := uc.CustomerRepository.Search(ctx, uc.DB, pagination)
	if err != nil {
		return nil, apperror.BadRequest(err)
	}
	return c, nil
}
func (uc *CustomerUseCase) Update(ctx context.Context, request *dtos.CreateOrUpdateCustomerRequest, id string) (*entities.Customer, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	c := new(entities.Customer)
	if err := uc.CustomerRepository.FindById(ctx, tx, c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Customer{}, apperror.BadRequest(apperror.ErrCustomerNotFound)
		}
		return &entities.Customer{}, apperror.InternalServerError(err)
	}

	c.Email = request.Email
	c.Name = request.Name

	if err := uc.CustomerRepository.Update(tx, c); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &entities.Customer{}, apperror.BadRequest(apperror.ErrCustomerEmailExist)
		}
		return &entities.Customer{}, apperror.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return &entities.Customer{}, apperror.BadRequest(err)
	}

	return converter.CustomerToResponse(c), nil
}
func (uc *CustomerUseCase) Delete(ctx context.Context, id string) error {
	c := new(entities.Customer)
	if err := uc.CustomerRepository.FindById(ctx, uc.DB, c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.BadRequest(apperror.ErrCustomerNotFound)
		}
		return apperror.InternalServerError(err)
	}

	if err := uc.CustomerRepository.Delete(uc.DB, c); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return apperror.BadRequest(apperror.ErrCustomerExistInvoice)
		}
		return apperror.BadRequest(err)
	}
	return nil
}
