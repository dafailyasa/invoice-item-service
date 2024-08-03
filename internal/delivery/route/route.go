package route

import (
	"fmt"

	"github.com/dafailyasa/invoice-item-service/internal/delivery"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	App         *echo.Echo
	CustomerHdl *delivery.CustomerHandler
	ItemHdl     *delivery.ItemHandler
	InvoiceHdl  *delivery.InvoiceHandler
}

func (r *Routes) Setup() {
	v := r.App.Group("api/v1")
	r.setupRoutes(v)
}

func (r *Routes) setupRoutes(v *echo.Group) {
	v.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, fmt.Sprint("PONG!! 👋"))
	})

	// start customers route
	customer := v.Group("/customers")
	customer.POST("", r.CustomerHdl.Create)
	customer.GET("/search", r.CustomerHdl.List)
	customer.GET("/:customerId", r.CustomerHdl.Detail)
	customer.PATCH("/:customerId", r.CustomerHdl.Update)
	customer.DELETE("/:customerId", r.CustomerHdl.Delete)
	// end customers route

	// start items route
	item := v.Group("/items")
	item.POST("", r.ItemHdl.Create)
	item.GET("/search", r.ItemHdl.List)
	item.PATCH("/:itemId", r.ItemHdl.Update)
	// end items route

	// start invoices route
	invoice := v.Group("/invoices")
	invoice.GET("/search", r.InvoiceHdl.List)
	invoice.POST("", r.InvoiceHdl.Create)
	invoice.PATCH("/:invoiceId", r.InvoiceHdl.Update)

}
