package entities

import (
	"time"

	"github.com/dafailyasa/invoice-item-service/pkg/constants"
)

type Customer struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:id" json:"id,omitempty"`
	Name      string     `gorm:"type:varchar(100);column:name" json:"name,omitempty"`
	Address   string     `gorm:"type:varchar(100);column:address" json:"address,omitempty"`
	Email     string     `gorm:"type:varchar(100);column:email" json:"email,omitempty"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt,omitempty"`
}

func (Customer) TableName() string {
	return constants.TableNameCustomers
}
