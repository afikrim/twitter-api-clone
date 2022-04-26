package repositories

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type TodoRepository interface {
	Create(ctx context.Context, dto *domains.CreateTodoDto) (*domains.Todo, error)
	FindAll(ctx context.Context) ([]domains.Todo, error)
	FindOne(ctx context.Context, id uint64) (*domains.Todo, error)
	Update(ctx context.Context, id uint64, dto *domains.UpdateTodoDto) (*domains.Todo, error)
	Remove(ctx context.Context, id uint64) error
}
