package dtos

import (
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type InvoiceItem struct {
	ID     int     `json:"id"`
	ItemID int     `json:"itemId"`
	Qty    int     `json:"qty"`
	Price  float64 `json:"price"`
}

func (ii InvoiceItem) Validate() error {
	return validation.ValidateStruct(&ii,
		validation.Field(&ii.ID, is.Int),
		validation.Field(&ii.ItemID, validation.Required),
		validation.Field(&ii.Qty, validation.Required),
		validation.Field(&ii.Price, validation.Required),
	)
}

func (ii InvoiceItem) MapToEntityInvoiceItem(invoiceId uint64) entities.InvoiceItem {
	return entities.InvoiceItem{
		ID:        uint64(ii.ID),
		ItemID:    uint64(ii.ItemID),
		InvoiceID: invoiceId,
		Amount:    utils.CalculateAmountAndQuantity(ii.Qty, ii.Price),
		Quantity:  ii.Qty,
		Price:     ii.Price,
	}
}
