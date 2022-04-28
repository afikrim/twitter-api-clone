package userfollowing_service

import (
	"context"
	"errors"
	"strconv"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/repositories"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
)

var (
	defaultLimit             = int(10)
	defaultOffset            = int(0)
	ErrInvalidUserFollowing  = errors.New("invalid user following")
	ErrUserFollowingNotFound = errors.New("user following not found")
	ErrUserNotFound          = errors.New("user not found")
)

type service struct {
	repo     repositories.UserFollowingRepository
	userRepo repositories.UserRepository
}

func NewUserFollowingService(repo repositories.UserFollowingRepository, userRepo repositories.UserRepository) *service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *service) Create(ctx context.Context, currentUserID string, followUserID string) error {
	if currentUserID == followUserID {
		return ErrInvalidUserFollowing
	}

	parsedCurrentUserID, err := strconv.ParseUint(currentUserID, 10, 64)
	if err != nil {
		return err
	}
	parsedFollowUserID, err := strconv.ParseUint(followUserID, 10, 64)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ctx, parsedFollowUserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserFollowingNotFound
	}

	err = s.repo.Create(ctx, parsedCurrentUserID, parsedFollowUserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FindAllFollowing(ctx context.Context, username string, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	if query.Limit == nil {
		query.Limit = &defaultLimit
	}
	if query.Offset == nil && query.Page != nil {
		offset := int(*query.Limit) * (int(*query.Page) - 1)
		query.Offset = &offset
	}
	if query.Offset == nil {
		query.Offset = &defaultOffset
	}

	users, cursor, err := s.repo.FindAllFollowing(ctx, user.ID, query)
	if err != nil {
		return nil, nil, err
	}

	return users, cursor, nil
}

func (s *service) FindAllFollowers(ctx context.Context, username string, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	if query.Limit == nil {
		query.Limit = &defaultLimit
	}
	if query.Offset == nil && query.Page != nil {
		offset := int(*query.Limit) * (int(*query.Page) - 1)
		query.Offset = &offset
	}
	if query.Offset == nil {
		query.Offset = &defaultOffset
	}

	users, cursor, err := s.repo.FindAllFollowers(ctx, user.ID, query)
	if err != nil {
		return nil, nil, err
	}

	return users, cursor, nil
}

func (s *service) Remove(ctx context.Context, currentUserID string, followUserID string) error {
	parsedCurrentUserID, err := strconv.ParseUint(currentUserID, 10, 64)
	if err != nil {
		return err
	}
	parsedFollowUserID, err := strconv.ParseUint(followUserID, 10, 64)
	if err != nil {
		return err
	}

	err = s.repo.Remove(ctx, parsedCurrentUserID, parsedFollowUserID)
	if err != nil {
		return err
	}

	return nil
}
