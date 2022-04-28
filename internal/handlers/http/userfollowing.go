package http_handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/services"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserFollowingHandler struct {
	service services.UserFollowingService
}

func NewUserFollowingHandler(service services.UserFollowingService) *UserFollowingHandler {
	return &UserFollowingHandler{
		service: service,
	}
}

func (h *UserFollowingHandler) Create(e echo.Context) error {
	ctx := context.Background()

	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(*domains.JwtCustomClaims)

	currentUserID := fmt.Sprint(claims.Session.UserID)
	followUserID := e.Param("credential")

	if err := h.service.Create(ctx, currentUserID, followUserID); err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully create user following"})
}

func (h *UserFollowingHandler) FindAllFollowing(e echo.Context) error {
	ctx := context.Background()

	username := e.Param("credential")
	query := &domains.QueryParamFollowDto{}
	if e.QueryParam("limit") != "" {
		limit, _ := strconv.Atoi(e.QueryParam("limit"))
		query.QueryParamPaginationDto.Limit = &limit
	}
	if e.QueryParam("offset") != "" {
		offset, _ := strconv.Atoi(e.QueryParam("offset"))
		query.QueryParamPaginationDto.Offset = &offset
	}
	if e.QueryParam("page") != "" {
		page, _ := strconv.Atoi(e.QueryParam("page"))
		query.QueryParamPaginationDto.Page = &page
	}

	users, cursor, err := h.service.FindAllFollowing(ctx, username, query)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully find all following", Data: map[string]interface{}{"users": users}, Meta: map[string]interface{}{"cursor": cursor}})
}

func (h *UserFollowingHandler) FindAllFollowers(e echo.Context) error {
	ctx := context.Background()

	username := e.Param("credential")
	query := &domains.QueryParamFollowDto{}
	if e.QueryParam("limit") != "" {
		limit, _ := strconv.Atoi(e.QueryParam("limit"))
		query.QueryParamPaginationDto.Limit = &limit
	}
	if e.QueryParam("offset") != "" {
		offset, _ := strconv.Atoi(e.QueryParam("offset"))
		query.QueryParamPaginationDto.Offset = &offset
	}
	if e.QueryParam("page") != "" {
		page, _ := strconv.Atoi(e.QueryParam("page"))
		query.QueryParamPaginationDto.Page = &page
	}

	users, cursor, err := h.service.FindAllFollowers(ctx, username, query)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully find all followers", Data: map[string]interface{}{"users": users}, Meta: map[string]interface{}{"cursor": cursor}})
}

func (h *UserFollowingHandler) Remove(e echo.Context) error {
	ctx := context.Background()

	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(*domains.JwtCustomClaims)

	currentUserID := fmt.Sprint(claims.Session.UserID)
	followUserID := e.Param("credential")

	if err := h.service.Remove(ctx, currentUserID, followUserID); err != nil {
		return e.JSON(http.StatusInternalServerError, &Response{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return e.JSON(http.StatusOK, &Response{Status: http.StatusOK, Message: "Successfully remove user following"})
}

func (h *UserFollowingHandler) RegisterRoutes(e *echo.Group) {
	group := e.Group("/users/:credential")

	group.POST("/follow", h.Create, IsLoggedIn)
	group.GET("/following", h.FindAllFollowing)
	group.GET("/followers", h.FindAllFollowers)
	group.POST("/unfollow", h.Remove, IsLoggedIn)
}
