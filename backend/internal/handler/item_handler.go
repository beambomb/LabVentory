package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"labventroy/backend/internal/apperror"
	"labventroy/backend/internal/model"
	"labventroy/backend/internal/service"
)

// ItemHandler handles HTTP CRUD operations for laboratory components.
type ItemHandler struct {
	svc *service.InventoryService
}

// NewItemHandler constructs an ItemHandler with its InventoryService dependency.
func NewItemHandler(svc *service.InventoryService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

// CreateItem handles POST /items.
func (h *ItemHandler) CreateItem(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var input model.CreateComponentInput
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return Error(apperror.BadRequest("invalid JSON request body")), nil
	}

	item, err := h.svc.CreateComponent(ctx, input)
	if err != nil {
		return Error(err), nil
	}
	return JSON(http.StatusCreated, item), nil
}

// ListItems handles GET /items.
func (h *ItemHandler) ListItems(ctx context.Context, _ events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	items, err := h.svc.ListComponents(ctx)
	if err != nil {
		return Error(err), nil
	}
	return JSON(http.StatusOK, map[string]any{
		"items": items,
		"count": len(items),
	}), nil
}

// GetItem handles GET /items/{id}.
func (h *ItemHandler) GetItem(ctx context.Context, req events.APIGatewayV2HTTPRequest, id string) (events.APIGatewayV2HTTPResponse, error) {
	item, err := h.svc.GetComponent(ctx, id)
	if err != nil {
		return Error(err), nil
	}
	return JSON(http.StatusOK, item), nil
}

// UpdateItem handles PUT /items/{id}.
func (h *ItemHandler) UpdateItem(ctx context.Context, req events.APIGatewayV2HTTPRequest, id string) (events.APIGatewayV2HTTPResponse, error) {
	var input model.UpdateComponentInput
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return Error(apperror.BadRequest("invalid JSON request body")), nil
	}

	item, err := h.svc.UpdateComponent(ctx, id, input)
	if err != nil {
		return Error(err), nil
	}
	return JSON(http.StatusOK, item), nil
}

// DeleteItem handles DELETE /items/{id}.
func (h *ItemHandler) DeleteItem(ctx context.Context, _ events.APIGatewayV2HTTPRequest, id string) (events.APIGatewayV2HTTPResponse, error) {
	if err := h.svc.DeleteComponent(ctx, id); err != nil {
		return Error(err), nil
	}
	return JSON(http.StatusOK, map[string]string{
		"message": "component deleted",
		"item_id": id,
	}), nil
}

// CreateImageUploadURL handles POST /items/{id}/upload-url.
func (h *ItemHandler) CreateImageUploadURL(ctx context.Context, req events.APIGatewayV2HTTPRequest, id string) (events.APIGatewayV2HTTPResponse, error) {
	var input model.UploadURLInput
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return Error(apperror.BadRequest("invalid JSON request body")), nil
	}

	out, err := h.svc.RequestImageUploadURL(ctx, id, input)
	if err != nil {
		return Error(err), nil
	}
	return JSON(http.StatusOK, out), nil
}
