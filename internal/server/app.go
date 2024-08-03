package server

import (
	"github.com/dafailyasa/invoice-item-service/internal/delivery"
	"github.com/dafailyasa/invoice-item-service/internal/delivery/route"
	"github.com/dafailyasa/invoice-item-service/internal/repository"
	elastic_repository "github.com/dafailyasa/invoice-item-service/internal/repository/es"
	"github.com/dafailyasa/invoice-item-service/internal/usecase"
)

func (s *Server) MapApps() {
	// repositories
	customerRepo := repository.NewCustomerRepository(s.Logger)
	itemRepo := repository.NewItemRepository(s.Logger)
	invoiceRepo := repository.NewInvoiceRepository(s.Logger)
	invoiceItemRepo := repository.NewInvoiceItemRepository(s.Logger)
	elasticRepo := elastic_repository.NewElasticSearchRepository(s.ES, &s.Cfg.ElasticSearch.Index)

	// use cases
	customerUC := usecase.NewCustomerUseCase(s.DB, customerRepo)
	itemUC := usecase.NewItemUseCase(s.DB, itemRepo)
	invoiceItemUC := usecase.NewInvoiceItemUseCase(s.DB, invoiceItemRepo)
	invoiceUC := usecase.NewInvoiceUseCase(s.DB, customerRepo, itemRepo, invoiceRepo, elasticRepo, invoiceItemUC, itemUC)

	// handlers
	customerHdl := delivery.NewCustomerHandler(customerUC)
	itemHdl := delivery.NewItemHandler(itemUC)
	invoiceHdl := delivery.NewInvoiceHandler(invoiceUC)

	routes := route.Routes{
		App:         s.Echo,
		CustomerHdl: customerHdl,
		ItemHdl:     itemHdl,
		InvoiceHdl:  invoiceHdl,
	}

	routes.Setup()
}
