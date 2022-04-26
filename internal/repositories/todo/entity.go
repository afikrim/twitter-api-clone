package todo_repository

import (
	"time"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type Todo struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Title     string     `gorm:"column:title;not null"`
	Completed bool       `gorm:"column:completed;not null;default:false"`
	CreatedAt *time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (t Todo) TableName() string {
	return "todos"
}

func (t *Todo) ToDomain() *domains.Todo {
	return &domains.Todo{
		ID:        t.ID,
		Title:     t.Title,
		Completed: t.Completed,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (Todo) FromDomain(d *domains.Todo) *Todo {
	return &Todo{
		ID:        d.ID,
		Title:     d.Title,
		Completed: d.Completed,
	}
}

func (Todo) FromDomainWithTimestamps(d *domains.Todo) *Todo {
	parsedCreatedAt, err := time.ParseInLocation("2006-01-02 15:04:05", d.CreatedAt, time.Local)
	if err != nil {
		panic(err)
	}

	parsedUpdatedAt, err := time.ParseInLocation("2006-01-02 15:04:05", d.UpdatedAt, time.Local)
	if err != nil {
		panic(err)
	}

	return &Todo{
		ID:        d.ID,
		Title:     d.Title,
		Completed: d.Completed,
		CreatedAt: &parsedCreatedAt,
		UpdatedAt: &parsedUpdatedAt,
	}
}
