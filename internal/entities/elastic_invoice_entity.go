package entities

import "github.com/dafailyasa/invoice-item-service/pkg/constants"

type ElasticInvoice struct {
	ID          uint64  `json:"id"`
	InvoiceID   string  `json:"invoiceId"`
	IssueDate   string  `json:"issueDate"`
	DueDate     string  `json:"dueDate"`
	Subject     string  `json:"subject"`
	ItemCount   int     `json:"itemCount"`
	Customer    string  `json:"customer"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"totalAmount"`
}

func MapToElasticInvoiceEntity(i *Invoice) *ElasticInvoice {
	return &ElasticInvoice{
		ID:          i.ID,
		InvoiceID:   i.InvoiceID,
		IssueDate:   i.IssueDate.Format(constants.DateFormat),
		DueDate:     i.DueDate.Format(constants.DateFormat),
		ItemCount:   i.ItemCount,
		Subject:     i.Subject,
		Customer:    i.Customer.Name,
		Status:      i.Status,
		TotalAmount: i.TotalAmount,
	}
}
