package todo_service

import (
	"context"
	"strconv"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/repositories"
)

type service struct {
	repo repositories.TodoRepository
}

func NewTodoService(repo repositories.TodoRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, dto *domains.CreateTodoDto) (*domains.Todo, error) {
	todo, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *service) FindAll(ctx context.Context) ([]domains.Todo, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Update(ctx context.Context, id string, dto *domains.UpdateTodoDto) (*domains.Todo, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	todo, err := s.repo.Update(ctx, parsedId, dto)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *service) Remove(ctx context.Context, id string) error {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	return s.repo.Remove(ctx, parsedId)
}
