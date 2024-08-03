package dtos

import (
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateOrUpdateCustomerRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

func (c CreateOrUpdateCustomerRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Address, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}

func (c CreateOrUpdateCustomerRequest) MapToEntityCustomer() entities.Customer {
	return entities.Customer{
		Name:    c.Name,
		Email:   c.Email,
		Address: c.Address,
	}
}
