package delivery

import (
	dtos "github.com/dafailyasa/invoice-item-service/internal/dtos/item"
	"github.com/dafailyasa/invoice-item-service/internal/usecase"
	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	"github.com/dafailyasa/invoice-item-service/pkg/response"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	UseCase *usecase.ItemUseCase
}

func NewItemHandler(useCase *usecase.ItemUseCase) *ItemHandler {
	return &ItemHandler{
		UseCase: useCase,
	}
}

func (h *ItemHandler) Create(c echo.Context) error {
	request := new(dtos.CreateOrUpdateItemRequest)
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

func (h *ItemHandler) List(c echo.Context) error {
	var p pagination.PaginationRequest
	if err := c.Bind(&p); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	data, err := h.UseCase.List(c.Request().Context(), &p)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(data, p).Send(c)
}

func (h *ItemHandler) Update(c echo.Context) error {
	request := new(dtos.CreateOrUpdateItemRequest)
	if err := c.Bind(request); err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	if err := request.Validate(); err != nil {
		return response.ErrorBuilder(apperror.UnprocessableEntity(err)).Send(c)
	}

	id := c.Param("itemId")
	data, err := h.UseCase.Update(c.Request().Context(), request, id)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(data).Send(c)
}
