package converter

import (
	"github.com/dafailyasa/invoice-item-service/internal/entities"
)

func ItemToResponse(i *entities.Item) *entities.Item {
	return &entities.Item{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

//
//func CustomersToResponse(customers *[]entities.Customer) (response []entities.Customer) {
//	for _, c := range *customers {
//		response = append(response, *CustomerToResponse(&c))
//	}
//
//	return response
//}
