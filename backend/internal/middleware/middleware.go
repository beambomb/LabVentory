package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// HandlerFunc represents an API Gateway v2 Lambda handler function.
type HandlerFunc func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

// Middleware wraps a HandlerFunc with cross-cutting behavior.
type Middleware func(HandlerFunc) HandlerFunc

// Chain applies middlewares in order so the first middleware is the outermost wrapper.
func Chain(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// CORSHeaders returns standard CORS headers for the React SPA client.
func CORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
	}
}

// WithCORS handles OPTIONS preflight requests and ensures CORS headers on all responses.
func WithCORS() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			if req.RequestContext.HTTP.Method == http.MethodOptions {
				return events.APIGatewayV2HTTPResponse{
					StatusCode: http.StatusNoContent,
					Headers:    CORSHeaders(),
				}, nil
			}

			resp, err := next(ctx, req)
			if resp.Headers == nil {
				resp.Headers = make(map[string]string)
			}
			for k, v := range CORSHeaders() {
				if _, exists := resp.Headers[k]; !exists {
					resp.Headers[k] = v
				}
			}
			return resp, err
		}
	}
}

// WithLogger logs incoming requests, status codes, and latency to CloudWatch Logs / stdout.
func WithLogger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			log.Printf("method=%s path=%s status=%d duration=%s",
				req.RequestContext.HTTP.Method,
				req.RawPath,
				resp.StatusCode,
				time.Since(start).Round(time.Millisecond),
			)
			return resp, err
		}
	}
}

// WithRecovery catches panics and returns a 500 JSON response instead of crashing.
func WithRecovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (resp events.APIGatewayV2HTTPResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic recovered: %v\n%s", r, string(debug.Stack()))
					body, _ := json.Marshal(map[string]string{
						"code":  "INTERNAL_ERROR",
						"error": "internal server error",
					})
					resp = events.APIGatewayV2HTTPResponse{
						StatusCode: http.StatusInternalServerError,
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						Body: string(body),
					}
					err = nil
				}
			}()
			return next(ctx, req)
		}
	}
}
