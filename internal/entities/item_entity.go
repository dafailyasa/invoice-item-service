package entities

import (
	"time"

	"github.com/dafailyasa/invoice-item-service/pkg/constants"
)

type Item struct {
	ID          uint64     `gorm:"primary_key;auto_increment;column:id" json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(100);column:name" json:"name,omitempty"`
	Description string     `gorm:"type:text;column:name" json:"description,omitempty"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updatedAt" json:"updatedAt,omitempty"`
}

func (Item) TableName() string {
	return constants.TableItems
}
