package main

import (
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"labventroy/backend/internal/config"
	"labventroy/backend/internal/handler"
	"labventroy/backend/internal/repository"
	"labventroy/backend/internal/route"
	"labventroy/backend/internal/service"
)

func main() {
	cfg := config.Load()

	repo := repository.NewMemoryRepository(true)
	svc := service.NewInventoryService(repo)
	itemHandler := handler.NewItemHandler(svc)
	router := route.NewRouter(itemHandler)

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		lambdaReq := events.APIGatewayV2HTTPRequest{
			RawPath: r.URL.Path,
			Body:    string(bodyBytes),
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: r.Method,
					Path:   r.URL.Path,
				},
			},
		}

		resp, err := router(r.Context(), lambdaReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write([]byte(resp.Body))
	})

	addr := ":" + cfg.Port
	log.Printf("LabComponentInventory local demo server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
