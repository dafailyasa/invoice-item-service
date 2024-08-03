package dtos

import (
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateOrUpdateItemRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=10000"`
}

func (i CreateOrUpdateItemRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Description, validation.Required),
	)
}

func (i CreateOrUpdateItemRequest) MapToEntityItem() entities.Item {
	return entities.Item{
		Name:        i.Name,
		Description: i.Description,
	}
}
