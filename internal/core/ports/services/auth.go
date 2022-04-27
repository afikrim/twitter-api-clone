package services

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type AuthService interface {
	Register(ctx context.Context, dto *domains.RegisterDto) error
	Login(ctx context.Context, dto *domains.LoginDto) (*domains.AuthWithRefresh, error)
	Refresh(ctx context.Context, refreshToken string) (*domains.AuthWithoutRefresh, error)
	Logout(ctx context.Context, refreshToken string) error
}
