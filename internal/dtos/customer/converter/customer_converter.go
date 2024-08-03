package converter

import (
	"github.com/dafailyasa/invoice-item-service/internal/entities"
)

func CustomerToResponse(c *entities.Customer) *entities.Customer {
	return &entities.Customer{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
