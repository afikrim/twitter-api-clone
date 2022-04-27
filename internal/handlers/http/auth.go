package http_handler

import (
	"context"
	"net/http"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(e echo.Context) error {
	ctx := context.Background()

	dto := new(domains.RegisterDto)
	if err := e.Bind(dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	if err := h.service.Register(ctx, dto); err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, &Response{Status: http.StatusCreated, Message: "Successfully register user"})
}

func (h *AuthHandler) Login(e echo.Context) error {
	ctx := context.Background()

	dto := new(domains.LoginDto)
	if err := e.Bind(dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	auth, err := h.service.Login(ctx, dto)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully login user", Data: auth})
}

func (h *AuthHandler) Refresh(e echo.Context) error {
	ctx := context.Background()

	refreshToken := e.Get("refresh_token").(map[string]interface{})["refresh_token"].(string)
	auth, err := h.service.Refresh(ctx, refreshToken)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully refresh token", Data: auth})
}

func (h *AuthHandler) Logout(e echo.Context) error {
	ctx := context.Background()

	refreshToken := e.Get("refresh_token").(map[string]interface{})["refresh_token"].(string)
	if err := h.service.Logout(ctx, refreshToken); err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully logout user"})
}

func (h *AuthHandler) RegisterRoutes(e *echo.Group) {
	group := e.Group("/auth")

	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/refresh", h.Refresh, ValidateRefreshToken)
	group.POST("/logout", h.Logout, ValidateRefreshToken)
}
