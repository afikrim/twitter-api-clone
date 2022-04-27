package repositories

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type SessionRepository interface {
	Create(ctx context.Context, refreshToken string, session *domains.Session) error
	FindByRefreshToken(ctx context.Context, refreshToken string) (*domains.Session, error)
	Remove(ctx context.Context, refreshToken string) error
}
