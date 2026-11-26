package repository

import (
	"context"
	"sort"
	"sync"

	"labventroy/backend/internal/apperror"
	"labventroy/backend/internal/model"
)

// MemoryRepository is a thread-safe in-memory ComponentRepository for local demos and tests.
type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]model.Component
}

// NewMemoryRepository creates an in-memory repository, optionally seeded with sample lab items.
func NewMemoryRepository(seed bool) *MemoryRepository {
	repo := &MemoryRepository{
		items: make(map[string]model.Component),
	}
	if seed {
		samples := []model.Component{
			{
				ItemID:         "CMP-001",
				Name:           "STM32F401 Blackpill",
				Category:       model.CategoryMicrocontroller,
				StockQty:       15,
				MinStock:       5,
				DrawerLocation: "RAK-A2",
				Condition:      model.ConditionGood,
				UpdatedAt:      "2026-11-25T23:20:00Z",
			},
			{
				ItemID:         "CMP-002",
				Name:           "ESP32-S3 DevKitC-1",
				Category:       model.CategoryMicrocontroller,
				StockQty:       3,
				MinStock:       5,
				DrawerLocation: "RAK-A1",
				Condition:      model.ConditionNew,
				UpdatedAt:      "2026-11-25T23:20:00Z",
			},
			{
				ItemID:         "CMP-003",
				Name:           "DHT22 Temperature & Humidity",
				Category:       model.CategorySensor,
				StockQty:       12,
				MinStock:       4,
				DrawerLocation: "RAK-B3",
				Condition:      model.ConditionGood,
				UpdatedAt:      "2026-11-25T23:20:00Z",
			},
			{
				ItemID:         "CMP-004",
				Name:           "Jumper Wire Male-Female 20cm",
				Category:       model.CategoryCable,
				StockQty:       4,
				MinStock:       10,
				DrawerLocation: "RAK-C1",
				Condition:      model.ConditionGood,
				UpdatedAt:      "2026-11-25T23:20:00Z",
			},
			{
				ItemID:         "CMP-005",
				Name:           "Resistor 10k Ohm 1/4W (Pack)",
				Category:       model.CategoryPassive,
				StockQty:       25,
				MinStock:       10,
				DrawerLocation: "RAK-D4",
				Condition:      model.ConditionNew,
				UpdatedAt:      "2026-11-25T23:20:00Z",
			},
		}
		for _, item := range samples {
			repo.items[item.ItemID] = item
		}
	}
	return repo
}

func (r *MemoryRepository) Create(_ context.Context, item model.Component) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ItemID]; exists {
		return apperror.Conflict("component with item_id " + item.ItemID + " already exists")
	}
	r.items[item.ItemID] = item
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]model.Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Component, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ItemID < result[j].ItemID
	})
	return result, nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id string) (*model.Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, apperror.NotFound("component " + id + " not found")
	}
	copyItem := item
	return &copyItem, nil
}

func (r *MemoryRepository) Update(_ context.Context, id string, input model.UpdateComponentInput, updatedAt string) (*model.Component, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.items[id]
	if !exists {
		return nil, apperror.NotFound("component " + id + " not found")
	}

	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Category != nil {
		item.Category = *input.Category
	}
	if input.StockQty != nil {
		item.StockQty = *input.StockQty
	}
	if input.MinStock != nil {
		item.MinStock = *input.MinStock
	}
	if input.DrawerLocation != nil {
		item.DrawerLocation = *input.DrawerLocation
	}
	if input.Condition != nil {
		item.Condition = *input.Condition
	}
	if input.ImageKey != nil {
		item.ImageKey = *input.ImageKey
	}
	item.UpdatedAt = updatedAt

	r.items[id] = item
	copyItem := item
	return &copyItem, nil
}

func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return apperror.NotFound("component " + id + " not found")
	}
	delete(r.items, id)
	return nil
}
