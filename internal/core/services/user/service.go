package user_service

import (
	"context"
	"strconv"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/repositories"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
	"github.com/go-playground/validator/v10"
)

var (
	defaultLimit   = int(10)
	defaultOffset  = int(0)
	defaultSortBy  = "id"
	defaultOrderBy = "asc"
)

type service struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) FindAll(ctx context.Context, query *domains.QueryParamUserDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error) {
	validate := validator.New()
	if err := validate.Struct(query); err != nil {
		return nil, nil, err
	}

	if query.Limit == nil {
		query.Limit = &defaultLimit
	}

	if query.Offset == nil && query.Page != nil {
		offset := int(*query.Limit) * (int(*query.Page) - 1)
		query.Offset = &offset
	}

	if query.Offset == nil && query.Page == nil {
		query.Offset = &defaultOffset
	}

	if query.SortBy == "" {
		query.SortBy = defaultSortBy
	}

	if query.OrderBy == "" {
		query.OrderBy = defaultOrderBy
	}

	users, cursor, err := s.repo.FindAll(ctx, query)
	if err != nil {
		return nil, nil, err
	}

	return users, cursor, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*domains.User, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(ctx, parsedId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) FindByUsername(ctx context.Context, username string) (*domains.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, id string, dto *domains.UpdateUserDto) (*domains.User, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Update(ctx, parsedId, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateCredential(ctx context.Context, id string, dto *domains.UpdateUserCredentialDto) (*domains.User, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.UpdateCredential(ctx, parsedId, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdatePassword(ctx context.Context, id string, dto *domains.UpdateUserPasswordDto) (*domains.User, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.UpdatePassword(ctx, parsedId, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) SoftRemove(ctx context.Context, id string) error {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	err = s.repo.SoftRemove(ctx, parsedId)
	if err != nil {
		return err
	}

	return nil
}
