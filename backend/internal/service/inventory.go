package service

import (
	"context"
	"strings"
	"time"

	"labventroy/backend/internal/apperror"
	"labventroy/backend/internal/model"
	"labventroy/backend/internal/repository"
)

// InventoryService encapsulates business logic for laboratory component inventory.
type InventoryService struct {
	repo repository.ComponentRepository
	now  func() time.Time
}

// NewInventoryService creates a new InventoryService with the given repository.
func NewInventoryService(repo repository.ComponentRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *InventoryService) CreateComponent(ctx context.Context, input model.CreateComponentInput) (*model.Component, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	item := input.ToComponent(s.now().Format(time.RFC3339))
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *InventoryService) ListComponents(ctx context.Context) ([]model.Component, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []model.Component{}, nil
	}
	return items, nil
}

func (s *InventoryService) GetComponent(ctx context.Context, id string) (*model.Component, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, apperror.Validation("item_id path parameter is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *InventoryService) UpdateComponent(ctx context.Context, id string, input model.UpdateComponentInput) (*model.Component, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, apperror.Validation("item_id path parameter is required")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, input, s.now().Format(time.RFC3339))
}

func (s *InventoryService) DeleteComponent(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperror.Validation("item_id path parameter is required")
	}
	return s.repo.Delete(ctx, id)
}
