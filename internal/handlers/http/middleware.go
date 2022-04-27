package http_handler

import (
	"net/http"
	"strings"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningMethod: middleware.AlgorithmHS256,
	SigningKey:    []byte("secret"),
	TokenLookup:   "header:" + echo.HeaderAuthorization,
	AuthScheme:    "Bearer",
	Claims:        &domains.JwtCustomClaims{},
})

func ValidateRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		var refreshToken string

		authHeader := e.Request().Header.Get("Authorization")
		if authHeader == "" {
			refreshToken = e.QueryParam("refresh_token")
			if refreshToken == "" {
				return e.JSON(http.StatusBadRequest, &Response{Status: http.StatusBadRequest, Message: "Refresh token is required"})
			}
		}
		if refreshToken == "" && authHeader != "" {
			refreshToken = strings.Replace(authHeader, "Bearer ", "", 1)
		}

		e.Set("refresh_token", map[string]interface{}{"refresh_token": refreshToken})
		return next(e)
	}
}
