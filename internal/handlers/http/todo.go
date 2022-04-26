package http_handler

import (
	"context"
	"net/http"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/services"
	"github.com/labstack/echo"
)

type TodoHttpHandler struct {
	service services.TodoService
}

func NewTodoHttpHandler(service services.TodoService) *TodoHttpHandler {
	return &TodoHttpHandler{
		service: service,
	}
}

func (h *TodoHttpHandler) Create(e echo.Context) error {
	ctx := context.Background()

	dto := domains.CreateTodoDto{}
	if err := e.Bind(&dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	todo, err := h.service.Create(ctx, &dto)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, &Response{Status: http.StatusCreated, Message: "Todo created", Data: map[string]interface{}{"todo": todo}})
}

func (h *TodoHttpHandler) FindAll(e echo.Context) error {
	ctx := context.Background()

	todos, err := h.service.FindAll(ctx)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Todo retrieved", Data: map[string]interface{}{"todos": todos}})
}

func (h *TodoHttpHandler) Update(e echo.Context) error {
	ctx := context.Background()

	dto := domains.UpdateTodoDto{}
	if err := e.Bind(&dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id is required"})
	}

	todo, err := h.service.Update(ctx, id, &dto)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Todo updated", Data: map[string]interface{}{"todo": todo}})
}

func (h *TodoHttpHandler) Remove(e echo.Context) error {
	ctx := context.Background()

	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id is required"})
	}

	err := h.service.Remove(ctx, id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Todo removed"})
}
