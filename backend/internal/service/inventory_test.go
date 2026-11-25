package service

import (
	"context"
	"testing"
	"time"

	"labventroy/backend/internal/model"
	"labventroy/backend/internal/repository"
)

func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

func TestInventoryService_CRUD(t *testing.T) {
	repo := repository.NewMemoryRepository(false)
	svc := NewInventoryService(repo)
	fixedTime := time.Date(2026, 11, 25, 23, 20, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedTime }

	ctx := context.Background()

	// Create
	created, err := svc.CreateComponent(ctx, model.CreateComponentInput{
		ItemID:         "CMP-100",
		Name:           "ESP32-WROOM-32",
		Category:       model.CategoryMicrocontroller,
		StockQty:       intPtr(10),
		MinStock:       intPtr(3),
		DrawerLocation: "RAK-A1",
		Condition:      model.ConditionNew,
	})
	if err != nil {
		t.Fatalf("CreateComponent failed: %v", err)
	}
	if created.UpdatedAt != "2026-11-25T23:20:00Z" {
		t.Fatalf("expected updated_at 2026-11-25T23:20:00Z, got %s", created.UpdatedAt)
	}

	// Duplicate create should fail
	if _, err := svc.CreateComponent(ctx, model.CreateComponentInput{
		ItemID:         "CMP-100",
		Name:           "ESP32-WROOM-32",
		Category:       model.CategoryMicrocontroller,
		StockQty:       intPtr(10),
		MinStock:       intPtr(3),
		DrawerLocation: "RAK-A1",
		Condition:      model.ConditionNew,
	}); err == nil {
		t.Fatal("expected conflict error on duplicate item_id")
	}

	// List
	list, err := svc.ListComponents(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 item, got %d, err=%v", len(list), err)
	}

	// Update
	cond := model.ConditionGood
	updated, err := svc.UpdateComponent(ctx, "CMP-100", model.UpdateComponentInput{
		StockQty:       intPtr(2),
		DrawerLocation: strPtr("RAK-A9"),
		Condition:      &cond,
	})
	if err != nil {
		t.Fatalf("UpdateComponent failed: %v", err)
	}
	if updated.StockQty != 2 || updated.DrawerLocation != "RAK-A9" || updated.Condition != model.ConditionGood {
		t.Fatalf("unexpected updated fields: %+v", updated)
	}

	// Delete
	if err := svc.DeleteComponent(ctx, "CMP-100"); err != nil {
		t.Fatalf("DeleteComponent failed: %v", err)
	}
	if _, err := svc.GetComponent(ctx, "CMP-100"); err == nil {
		t.Fatal("expected not found after delete")
	}
}
