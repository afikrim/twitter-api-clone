package services

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
)

type UserFollowingService interface {
	Create(ctx context.Context, currentUserID string, followUserID string) error
	FindAllFollowing(ctx context.Context, username string, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error)
	FindAllFollowers(ctx context.Context, username string, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error)
	Remove(ctx context.Context, currentUserID string, followUserID string) error
}
