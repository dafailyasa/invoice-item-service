package dtos

import (
	"fmt"
	dtos "github.com/dafailyasa/invoice-item-service/internal/dtos/invoice-item"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/pkg/constants"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samber/lo"
	"time"
)

type CreateOrUpdateInvoiceRequest struct {
	CustomerID int                `json:"customerId"`
	Subject    string             `json:"subject"`
	Status     string             `json:"status"`
	DueDate    string             `json:"dueDate"`
	IssueDate  string             `json:"issueDate"`
	Items      []dtos.InvoiceItem `json:"items"`
}

func (i CreateOrUpdateInvoiceRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.CustomerID, validation.Required),
		validation.Field(&i.Subject, validation.Required),
		validation.Field(&i.Status, validation.Required, validation.In(constants.InvStatusUnpaid, constants.InvStatusPaid).
			Error(
				fmt.Sprintf("invalid status type. should be (%s, %s)",
					constants.InvStatusUnpaid,
					constants.InvStatusPaid,
				),
			),
		),
		validation.Field(&i.DueDate, validation.Required, validation.Date(constants.DateFormatSQL)),
		validation.Field(&i.IssueDate, validation.Required, validation.Date(constants.DateFormatSQL)),
		validation.Field(&i.Items,
			validation.Required,
			validation.Each(
				validation.By(func(val interface{}) error {
					item := val.(dtos.InvoiceItem)
					return item.Validate()
				}),
			),
		),
	)
}

func (i CreateOrUpdateInvoiceRequest) MapToEntityInvoice() entities.Invoice {
	invoice := fmt.Sprintf("INV-%d", time.Now().UnixNano())

	dueDate := lo.Must(time.Parse(constants.DateFormatSQL, i.DueDate))
	issueDate := lo.Must(time.Parse(constants.DateFormatSQL, i.IssueDate))

	return entities.Invoice{
		CustomerID: uint64(i.CustomerID),
		InvoiceID:  invoice,
		Subject:    i.Subject,
		Status:     i.Status,
		DueDate:    dueDate,
		IssueDate:  issueDate,
		ItemCount:  len(i.Items),
	}
}

func (i CreateOrUpdateInvoiceRequest) MapToEntityInvoiceUpdate(oldInv *entities.Invoice) entities.Invoice {

	dueDate := lo.Must(time.Parse(constants.DateFormatSQL, i.DueDate))
	issueDate := lo.Must(time.Parse(constants.DateFormatSQL, i.IssueDate))

	return entities.Invoice{
		ID:         oldInv.ID,
		InvoiceID:  oldInv.InvoiceID,
		CustomerID: uint64(i.CustomerID),
		Status:     oldInv.Status,
		Subject:    i.Subject,
		DueDate:    dueDate,
		IssueDate:  issueDate,
	}
}

func MapToInvoiceDetailResponse(i *entities.Invoice) *entities.Invoice {
	return &entities.Invoice{
		ID:           i.ID,
		InvoiceID:    i.InvoiceID,
		CustomerID:   i.CustomerID,
		Subject:      i.Subject,
		Status:       i.Status,
		TotalAmount:  i.TotalAmount,
		ItemCount:    i.ItemCount,
		DueDate:      i.DueDate,
		IssueDate:    i.IssueDate,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
		Customer:     i.Customer,
		InvoiceItems: i.InvoiceItems,
	}
}
