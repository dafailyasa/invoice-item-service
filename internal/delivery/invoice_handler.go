package delivery

import (
	dtos "github.com/dafailyasa/invoice-item-service/internal/dtos/invoice"
	"github.com/dafailyasa/invoice-item-service/internal/usecase"
	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	"github.com/dafailyasa/invoice-item-service/pkg/response"
	"github.com/labstack/echo/v4"
)

type InvoiceHandler struct {
	UseCase *usecase.InvoiceUseCase
}

func NewInvoiceHandler(useCase *usecase.InvoiceUseCase) *InvoiceHandler {
	return &InvoiceHandler{
		UseCase: useCase,
	}
}

func (h *InvoiceHandler) Create(c echo.Context) error {
	request := new(dtos.CreateOrUpdateInvoiceRequest)
	if err := c.Bind(request); err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	if err := request.Validate(); err != nil {
		return response.ErrorBuilder(apperror.UnprocessableEntity(err)).Send(c)
	}

	data, err := h.UseCase.Create(c.Request().Context(), request)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(data).Send(c)
}

func (h *InvoiceHandler) Update(c echo.Context) error {
	request := new(dtos.CreateOrUpdateInvoiceRequest)
	if err := c.Bind(request); err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	if err := request.Validate(); err != nil {
		return response.ErrorBuilder(apperror.UnprocessableEntity(err)).Send(c)
	}

	id := c.Param("invoiceId")
	data, err := h.UseCase.Update(c.Request().Context(), request, id)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(data).Send(c)
}

func (h *InvoiceHandler) List(c echo.Context) error {
	var p pagination.PaginationRequest
	if err := c.Bind(&p); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	data, err := h.UseCase.Search(c.Request().Context(), &p)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(data, p).Send(c)
}
