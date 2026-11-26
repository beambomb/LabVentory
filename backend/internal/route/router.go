package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"labventroy/backend/internal/apperror"
	"labventroy/backend/internal/handler"
	"labventroy/backend/internal/middleware"
)

// NewRouter builds the API Gateway v2 route dispatcher wrapped with standard middlewares.
func NewRouter(itemHandler *handler.ItemHandler) middleware.HandlerFunc {
	dispatch := func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		method := strings.ToUpper(req.RequestContext.HTTP.Method)
		path := normalizePath(req.RawPath, req.RequestContext.HTTP.Path)

		if path == "/items" {
			switch method {
			case http.MethodPost:
				return itemHandler.CreateItem(ctx, req)
			case http.MethodGet:
				return itemHandler.ListItems(ctx, req)
			default:
				return handler.Error(apperror.BadRequest("method not allowed for /items")), nil
			}
		}

		if strings.HasPrefix(path, "/items/") {
			rest := strings.TrimPrefix(path, "/items/")
			if strings.HasSuffix(rest, "/upload-url") {
				id := strings.TrimSuffix(rest, "/upload-url")
				if id == "" || strings.Contains(id, "/") {
					return handler.Error(apperror.NotFound("route not found")), nil
				}
				if method != http.MethodPost {
					return handler.Error(apperror.BadRequest("method not allowed for /items/{id}/upload-url")), nil
				}
				return itemHandler.CreateImageUploadURL(ctx, req, id)
			}

			id := rest
			if id == "" || strings.Contains(id, "/") {
				return handler.Error(apperror.NotFound("route not found")), nil
			}

			switch method {
			case http.MethodGet:
				return itemHandler.GetItem(ctx, req, id)
			case http.MethodPut:
				return itemHandler.UpdateItem(ctx, req, id)
			case http.MethodDelete:
				return itemHandler.DeleteItem(ctx, req, id)
			default:
				return handler.Error(apperror.BadRequest("method not allowed for /items/{id}")), nil
			}
		}

		return handler.Error(apperror.NotFound("route not found")), nil
	}

	return middleware.Chain(
		dispatch,
		middleware.WithRecovery(),
		middleware.WithLogger(),
		middleware.WithCORS(),
	)
}

func normalizePath(rawPath, httpPath string) string {
	p := rawPath
	if p == "" {
		p = httpPath
	}
	if len(p) > 1 {
		p = strings.TrimRight(p, "/")
	}
	return p
}
