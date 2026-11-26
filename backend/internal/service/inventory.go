package service

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"labventroy/backend/internal/apperror"
	"labventroy/backend/internal/model"
	"labventroy/backend/internal/repository"
	"labventroy/backend/internal/storage"
)

const defaultPresignTTL = 15 * time.Minute

// InventoryService encapsulates business logic for laboratory component inventory.
type InventoryService struct {
	repo    repository.ComponentRepository
	storage storage.MediaStorage
	now     func() time.Time
}

// NewInventoryService creates a new InventoryService with no-op media storage.
func NewInventoryService(repo repository.ComponentRepository) *InventoryService {
	return NewInventoryServiceWithStorage(repo, storage.NoopMediaStorage{})
}

// NewInventoryServiceWithStorage creates an InventoryService with both repository and media storage.
func NewInventoryServiceWithStorage(repo repository.ComponentRepository, media storage.MediaStorage) *InventoryService {
	if media == nil {
		media = storage.NoopMediaStorage{}
	}
	return &InventoryService{
		repo:    repo,
		storage: media,
		now:     func() time.Time { return time.Now().UTC() },
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
	for i := range items {
		s.attachImageURL(ctx, &items[i])
	}
	return items, nil
}

func (s *InventoryService) GetComponent(ctx context.Context, id string) (*model.Component, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, apperror.Validation("item_id path parameter is required")
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.attachImageURL(ctx, item)
	return item, nil
}

func (s *InventoryService) UpdateComponent(ctx context.Context, id string, input model.UpdateComponentInput) (*model.Component, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, apperror.Validation("item_id path parameter is required")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	item, err := s.repo.Update(ctx, id, input, s.now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	s.attachImageURL(ctx, item)
	return item, nil
}

func (s *InventoryService) DeleteComponent(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperror.Validation("item_id path parameter is required")
	}
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if existing.ImageKey != "" {
		_ = s.storage.DeleteObject(ctx, existing.ImageKey)
	}
	return nil
}

// RequestImageUploadURL verifies the component exists and generates a pre-signed S3 PUT URL.
func (s *InventoryService) RequestImageUploadURL(ctx context.Context, id string, input model.UploadURLInput) (*model.UploadURLOutput, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, apperror.Validation("item_id path parameter is required")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}

	cleanName := filepath.Base(input.Filename)
	key := "items/" + id + "/" + cleanName
	uploadURL, err := s.storage.GenerateUploadURL(ctx, key, input.ContentType, defaultPresignTTL)
	if err != nil {
		return nil, err
	}

	return &model.UploadURLOutput{
		UploadURL: uploadURL,
		ImageKey:  key,
		ExpiresIn: int(defaultPresignTTL.Seconds()),
	}, nil
}

func (s *InventoryService) attachImageURL(ctx context.Context, item *model.Component) {
	if item == nil || item.ImageKey == "" {
		return
	}
	if url, err := s.storage.GenerateDownloadURL(ctx, item.ImageKey, defaultPresignTTL); err == nil {
		item.ImageURL = url
	}
}
