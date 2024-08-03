package repository

import (
	"context"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	Repository[entities.Invoice]
	logger zerolog.Logger
}

func NewInvoiceRepository(logger zerolog.Logger) *InvoiceRepository {
	return &InvoiceRepository{
		logger: logger,
	}
}

func (r *InvoiceRepository) FindInvoiceWithRelations(ctx context.Context, db *gorm.DB, id uint64) (invoice entities.Invoice, err error) {
	err = db.WithContext(ctx).Where("id = ?", id).Preload("Customer").Preload("InvoiceItems.Item").First(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}
