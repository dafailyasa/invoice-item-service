package repository

import (
	"context"

	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/pkg/constants"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	"github.com/rs/zerolog"
	"github.com/sourcegraph/conc/pool"
	"gorm.io/gorm"
)

type ItemRepository struct {
	Repository[entities.Item]
	logger zerolog.Logger
}

func NewItemRepository(logger zerolog.Logger) *ItemRepository {
	return &ItemRepository{
		logger: logger,
	}
}

func (r *ItemRepository) Search(ctx context.Context, db *gorm.DB, pagination *pagination.PaginationRequest) (customers []entities.Item, err error) {
	workers := pool.New().WithContext(ctx)
	filter := r.filterSearch(pagination)

	workers.Go(func(ctx context.Context) error {
		return db.Table(constants.TableItems).WithContext(ctx).Scopes(filter).Order(pagination.GetSort()).Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&customers).Error
	})

	workers.Go(func(ctx context.Context) error {
		return db.Table(constants.TableItems).WithContext(ctx).Scopes(filter).Count(&pagination.Total).Error
	})

	if err := workers.Wait(); err != nil {
		return customers, err
	}

	return customers, nil
}

func (r *ItemRepository) filterSearch(pagination *pagination.PaginationRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if s := pagination.Keyword; s != "" {
			s = "%" + s + "%"
			tx = tx.Where("name LIKE ?", s)
		}
		return tx
	}
}
