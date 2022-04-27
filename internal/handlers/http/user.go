package http_handler

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/services"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) FindAll(e echo.Context) error {
	ctx := context.Background()

	query := &domains.QueryParamDto{
		OrderBy: e.QueryParam("orderby"),
		SortBy:  e.QueryParam("sortby"),
		Search:  e.QueryParam("search"),
	}
	if e.QueryParam("limit") != "" {
		limit, _ := strconv.Atoi(e.QueryParam("limit"))
		query.Limit = &limit
	}
	if e.QueryParam("offset") != "" {
		offset, _ := strconv.Atoi(e.QueryParam("offset"))
		query.Offset = &offset
	}
	if e.QueryParam("page") != "" {
		page, _ := strconv.Atoi(e.QueryParam("page"))
		query.Page = &page
	}

	users, err := h.service.FindAll(ctx, query)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully get all users", Data: map[string]interface{}{"users": users}})
}

func (h *UserHandler) FindByUsername(e echo.Context) error {
	ctx := context.Background()

	username := e.Param("credential")
	if username == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "username is required"})
	}

	user, err := h.service.FindByUsername(ctx, username)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully get user", Data: map[string]interface{}{"user": user.HidePassword()}})
}

func (h *UserHandler) Update(e echo.Context) error {
	ctx := context.Background()

	id := e.Param("credential")
	if id == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id is required"})
	}

	isNumberRegex := regexp.MustCompile(`^\d+$`)
	if !isNumberRegex.MatchString(id) {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id must be number"})
	}

	if err := h.ValidateUserOwner(id, e); err != nil {
		return err
	}

	dto := new(domains.UpdateUserDto)
	if err := e.Bind(dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user, err := h.service.Update(ctx, id, dto)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully update user", Data: map[string]interface{}{"user": user}})
}

func (h *UserHandler) UpdateCredential(e echo.Context) error {
	ctx := context.Background()

	id := e.Param("credential")
	if id == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id is required"})
	}

	isNumberRegex := regexp.MustCompile(`^\d+$`)
	if !isNumberRegex.MatchString(id) {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id must be number"})
	}

	if err := h.ValidateUserOwner(id, e); err != nil {
		return err
	}

	dto := new(domains.UpdateUserCredentialDto)
	if err := e.Bind(dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user, err := h.service.UpdateCredential(ctx, id, dto)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully update user credential", Data: map[string]interface{}{"user": user}})
}

func (h *UserHandler) UpdatePassword(e echo.Context) error {
	ctx := context.Background()

	id := e.Param("credential")
	if id == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id is required"})
	}

	isNumberRegex := regexp.MustCompile(`^\d+$`)
	if !isNumberRegex.MatchString(id) {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id must be number"})
	}

	if err := h.ValidateUserOwner(id, e); err != nil {
		return err
	}

	dto := new(domains.UpdateUserPasswordDto)
	if err := e.Bind(dto); err != nil {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user, err := h.service.UpdatePassword(ctx, id, dto)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully update user password", Data: map[string]interface{}{"user": user}})
}

func (h *UserHandler) SoftRemove(e echo.Context) error {
	ctx := context.Background()

	id := e.Param("credential")
	if id == "" {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id is required"})
	}

	isNumberRegex := regexp.MustCompile(`^\d+$`)
	if !isNumberRegex.MatchString(id) {
		return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "id must be number"})
	}

	if err := h.ValidateUserOwner(id, e); err != nil {
		return err
	}

	err := h.service.SoftRemove(ctx, id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully delete user"})
}

func (h *UserHandler) ValidateUserOwner(id string, e echo.Context) error {
	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(*domains.JwtCustomClaims)
	if fmt.Sprint(claims.Session.UserID) != id {
		return e.JSON(http.StatusForbidden, &Response{Status: http.StatusForbidden, Message: "You are not allowed to access this resource"})
	}

	return nil
}

func (h *UserHandler) RegisterRoutes(e *echo.Group) {
	group := e.Group("/users")

	group.GET("", h.FindAll)
	group.GET("/:credential", h.FindByUsername)
	group.PATCH("/:credential", h.Update, IsLoggedIn)
	group.PATCH("/:credential/credential", h.UpdateCredential, IsLoggedIn)
	group.PATCH("/:credential/password", h.UpdatePassword, IsLoggedIn)
	group.DELETE("/:credential", h.SoftRemove, IsLoggedIn)
}
