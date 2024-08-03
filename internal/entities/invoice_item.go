package entities

import (
	"time"

	"github.com/dafailyasa/invoice-item-service/pkg/constants"
)

type InvoiceItem struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:id"`
	InvoiceID uint64     `gorm:"type:bigint;column:invoiceId;default:null"`
	ItemID    uint64     `gorm:"type:bigint;column:itemId;default:null"`
	Quantity  int        `gorm:"type:bigint;column:quantity;default:null"`
	Price     float64    `gorm:"type:decimal(10,2);column:price"`
	Amount    float64    `gorm:"type:varchar(11);column:amount;default:null;generatedAlways"`
	CreatedAt time.Time  `gorm:"column:createdAt"`
	UpdatedAt *time.Time `gorm:"column:updatedAt"`

	Item Item `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"item,omitempty"`
}

func (InvoiceItem) TableName() string {
	return constants.TableInvoiceItems
}
