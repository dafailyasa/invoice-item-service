package usecase

import (
	"context"
	"github.com/dafailyasa/invoice-item-service/internal/repository"
	"gorm.io/gorm"
)

type InvoiceItemUseCase struct {
	DB                    *gorm.DB
	InvoiceItemRepository *repository.InvoiceItemRepository
}

func NewInvoiceItemUseCase(
	db *gorm.DB,
	invoiceItemRepository *repository.InvoiceItemRepository,
) *InvoiceItemUseCase {
	return &InvoiceItemUseCase{
		DB:                    db,
		InvoiceItemRepository: invoiceItemRepository,
	}
}

func (u *InvoiceItemUseCase) RemoveByIds(ctx context.Context, db *gorm.DB, invId uint64, ids []uint64) error {
	return u.InvoiceItemRepository.DeleteByIds(ctx, db, invId, ids)
}
