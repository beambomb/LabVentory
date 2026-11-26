package route

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"labventroy/backend/internal/handler"
	"labventroy/backend/internal/model"
	"labventroy/backend/internal/repository"
	"labventroy/backend/internal/service"
)

func setupRouter() func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	repo := repository.NewMemoryRepository(false)
	svc := service.NewInventoryService(repo)
	h := handler.NewItemHandler(svc)
	return NewRouter(h)
}

func makeReq(method, path, body string) events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest{
		RawPath: path,
		Body:    body,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: method,
				Path:   path,
			},
		},
	}
}

func TestRouter_CRUDEndpoints(t *testing.T) {
	router := setupRouter()
	ctx := context.Background()

	// 1. OPTIONS preflight returns 204 + CORS headers
	optResp, err := router(ctx, makeReq(http.MethodOptions, "/items", ""))
	if err != nil || optResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 for OPTIONS, got %d, err=%v", optResp.StatusCode, err)
	}
	if optResp.Headers["Access-Control-Allow-Origin"] != "*" {
		t.Fatalf("missing CORS header on OPTIONS")
	}

	// 2. POST /items -> 201 Created
	createBody := `{"item_id":"CMP-001","name":"STM32F401 Blackpill","category":"MICROCONTROLLER","stock_qty":15,"min_stock":5,"drawer_location":"RAK-A2","condition":"GOOD"}`
	createResp, err := router(ctx, makeReq(http.MethodPost, "/items", createBody))
	if err != nil || createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d (%s), err=%v", createResp.StatusCode, createResp.Body, err)
	}
	if createResp.Headers["Access-Control-Allow-Origin"] != "*" {
		t.Fatalf("missing CORS header on POST /items")
	}

	// 3. GET /items -> 200 OK
	listResp, err := router(ctx, makeReq(http.MethodGet, "/items", ""))
	if err != nil || listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", listResp.StatusCode)
	}
	var listPayload struct {
		Items []model.Component `json:"items"`
		Count int               `json:"count"`
	}
	if err := json.Unmarshal([]byte(listResp.Body), &listPayload); err != nil || listPayload.Count != 1 {
		t.Fatalf("expected count=1, got %+v, err=%v", listPayload, err)
	}

	// 4. GET /items/{id} -> 200 OK & 404 Not Found
	getResp, _ := router(ctx, makeReq(http.MethodGet, "/items/CMP-001", ""))
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", getResp.StatusCode)
	}
	notFoundResp, _ := router(ctx, makeReq(http.MethodGet, "/items/CMP-999", ""))
	if notFoundResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", notFoundResp.StatusCode)
	}

	// 4.5 POST /items/{id}/upload-url -> 200 OK & 404 Not Found
	uploadResp, _ := router(ctx, makeReq(http.MethodPost, "/items/CMP-001/upload-url", `{"filename":"stm32.png","content_type":"image/png"}`))
	if uploadResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on upload-url, got %d (%s)", uploadResp.StatusCode, uploadResp.Body)
	}
	var uploadOut model.UploadURLOutput
	if err := json.Unmarshal([]byte(uploadResp.Body), &uploadOut); err != nil || uploadOut.ImageKey != "items/CMP-001/stm32.png" {
		t.Fatalf("unexpected upload-url output: %+v, err=%v", uploadOut, err)
	}
	uploadMissingResp, _ := router(ctx, makeReq(http.MethodPost, "/items/CMP-999/upload-url", `{"filename":"stm32.png","content_type":"image/png"}`))
	if uploadMissingResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing item upload-url, got %d", uploadMissingResp.StatusCode)
	}

	// 5. PUT /items/{id} -> 200 OK
	updateResp, _ := router(ctx, makeReq(http.MethodPut, "/items/CMP-001", `{"stock_qty":3,"drawer_location":"RAK-B1","condition":"DEFECTIVE"}`))
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on PUT, got %d (%s)", updateResp.StatusCode, updateResp.Body)
	}

	// 6. DELETE /items/{id} -> 200 OK
	delResp, _ := router(ctx, makeReq(http.MethodDelete, "/items/CMP-001", ""))
	if delResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on DELETE, got %d", delResp.StatusCode)
	}

	// 7. GET after DELETE -> 404 Not Found
	afterDelResp, _ := router(ctx, makeReq(http.MethodGet, "/items/CMP-001", ""))
	if afterDelResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after DELETE, got %d", afterDelResp.StatusCode)
	}
}
