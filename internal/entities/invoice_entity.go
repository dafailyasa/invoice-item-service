package entities

import (
	"time"

	"github.com/dafailyasa/invoice-item-service/pkg/constants"
)

type Invoice struct {
	ID          uint64     `gorm:"primary_key;auto_increment;column:id" json:"id"`
	CustomerID  uint64     `gorm:"type:bigint;column:customerId" json:"customerId"`
	InvoiceID   string     `gorm:"type:varchar(200);column:invoiceId" json:"invoiceId"`
	Subject     string     `gorm:"type:text;column:subject" json:"subject"`
	Status      string     `gorm:"type:enum('Paid','Unpaid');column:status;default:'Unpaid'" json:"status"`
	TotalAmount float64    `gorm:"type:decimal(10,2); column:totalAmount" json:"totalAmount"`
	ItemCount   int        `gorm:"type:int;column:itemCount" json:"itemCount"`
	DueDate     time.Time  `gorm:"type:date;column:dueDate" json:"dueDate"`
	IssueDate   time.Time  `gorm:"type:date;column:issueDate" json:"issueDate"`
	CreatedAt   *time.Time `gorm:"column:createdAt" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updatedAt"`

	Customer     Customer      `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InvoiceItems []InvoiceItem `gorm:"foreignKey:InvoiceID;references:ID"`
}

func (Invoice) TableName() string {
	return constants.TableInvoices
}
