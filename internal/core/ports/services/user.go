package services

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type UserService interface {
	FindAll(ctx context.Context, query *domains.QueryParamDto) ([]domains.User, error)
	FindByID(ctx context.Context, id string) (*domains.User, error)
	FindByUsername(ctx context.Context, username string) (*domains.User, error)
	Update(ctx context.Context, id string, dto *domains.UpdateUserDto) (*domains.User, error)
	UpdateCredential(ctx context.Context, id string, dto *domains.UpdateUserCredentialDto) (*domains.User, error)
	UpdatePassword(ctx context.Context, id string, dto *domains.UpdateUserPasswordDto) (*domains.User, error)
	SoftRemove(ctx context.Context, id string) error
}
