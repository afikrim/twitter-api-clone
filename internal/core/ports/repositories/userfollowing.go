package repositories

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
)

type UserFollowingRepository interface {
	Create(ctx context.Context, currentUserID uint64, followUserID uint64) error
	FindAllFollowing(ctx context.Context, userID uint64, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error)
	FindAllFollowers(ctx context.Context, userID uint64, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error)
	Remove(ctx context.Context, currentUserID uint64, followUserID uint64) error
}
