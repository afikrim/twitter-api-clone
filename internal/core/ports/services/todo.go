package services

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type TodoService interface {
	Create(ctx context.Context, dto *domains.CreateTodoDto) (*domains.Todo, error)
	FindAll(ctx context.Context) ([]domains.Todo, error)
	Update(ctx context.Context, id string, dto *domains.UpdateTodoDto) (*domains.Todo, error)
	Remove(ctx context.Context, id string) error
}
