package http_handler

import (
	"context"
	"net/http"

	"github.com/afikrim/go-hexa-template/internal/core/ports/services"
	"github.com/labstack/echo/v4"
)

type handler struct {
	service services.CountryService
}

func NewCountryHandler(service services.CountryService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) FindAll(e echo.Context) error {
	ctx := context.Background()

	countries, err := h.service.FindAll(ctx)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully get all countries", Data: countries})
}

func (h *handler) RegisterRoutes(e *echo.Group) {
	countriesRoute := e.Group("/countries")

	countriesRoute.GET("", h.FindAll)
}
