package repositories

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type UserRepository interface {
	Create(ctx context.Context, dto *domains.RegisterDto) (*domains.User, error)
	FindAll(ctx context.Context, query *domains.QueryParamDto) ([]domains.User, error)
	FindByID(ctx context.Context, id int64) (*domains.User, error)
	FindByUsername(ctx context.Context, username string) (*domains.User, error)
	FindByEmail(ctx context.Context, email string) (*domains.User, error)
	FindByPhone(ctx context.Context, phone string) (*domains.User, error)
	Update(ctx context.Context, id int64, dto *domains.UpdateUserDto) (*domains.User, error)
	UpdateCredential(ctx context.Context, id int64, dto *domains.UpdateUserCredentialDto) (*domains.User, error)
	UpdatePassword(ctx context.Context, id int64, dto *domains.UpdateUserPasswordDto) (*domains.User, error)
	SoftRemove(ctx context.Context, id int64) error
}
