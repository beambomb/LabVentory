package repository

import (
	"context"

	"labventroy/backend/internal/model"
)

// ComponentRepository defines persistence operations for laboratory components.
type ComponentRepository interface {
	Create(ctx context.Context, item model.Component) error
	List(ctx context.Context) ([]model.Component, error)
	GetByID(ctx context.Context, id string) (*model.Component, error)
	Update(ctx context.Context, id string, input model.UpdateComponentInput, updatedAt string) (*model.Component, error)
	Delete(ctx context.Context, id string) error
}
