package repository

import (
	"context"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/pkg/constants"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type InvoiceItemRepository struct {
	Repository[entities.InvoiceItem]
	logger zerolog.Logger
}

func NewInvoiceItemRepository(logger zerolog.Logger) *InvoiceItemRepository {
	return &InvoiceItemRepository{
		logger: logger,
	}
}

func (repo *InvoiceItemRepository) DeleteByIds(ctx context.Context, db *gorm.DB, invId uint64, ids []uint64) error {
	return db.Table(constants.TableInvoiceItems).WithContext(ctx).Delete(entities.InvoiceItem{}, "invoiceId = ? AND id NOT IN (?)", invId, ids).Error
}
