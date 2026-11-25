# Progress Tracker — Backend (Go AWS Lambda + DynamoDB)

Target arsitektur: **Separation of Concerns (SoC)** untuk `LabComponentInventory` (Track C: Serverless Architecture).

## Checklist Implementasi

- [x] **Setup**: Struktur folder `backend/`, `go.mod`, `.gitignore`
- [ ] **Config** (`internal/config/config.go`): Loader environment variables (`TABLE_NAME`, `AWS_REGION`, `PORT`)
- [ ] **AppError** (`internal/apperror/error.go`): Standardized domain & API errors (`NotFound`, `Conflict`, `Validation`, `Internal`)
- [ ] **Model** (`internal/model/component.go` + `component_test.go`): Entity `Component`, enum `Category` & `Condition`, validasi `Create` & `Update`
- [ ] **Repository** (`internal/repository/repository.go`, `memory.go`, `dynamodb.go`): Interface `ComponentRepository`, in-memory seed store, dan implementasi AWS SDK for Go v2
- [ ] **Service** (`internal/service/inventory.go` + `inventory_test.go`): Business logic, timestamp ISO 8601 (`updated_at`), normalisasi input
- [ ] **Middleware** (`internal/middleware/middleware.go`): CORS headers, CloudWatch request logging, panic recovery
- [ ] **Handler** (`internal/handler/response.go`, `item_handler.go`): CRUD controllers & JSON response helpers
- [ ] **Route** (`internal/route/router.go` + `router_test.go`): API Gateway v2 route dispatcher (`POST/GET/PUT/DELETE /items`)
- [ ] **Entrypoints** (`cmd/lambda/main.go`, `cmd/local/main.go`): Lambda `provided.al2023` & local demo HTTP server
- [ ] **Verifikasi**: `go test -v ./...` lulus & cross-compile `GOOS=linux GOARCH=amd64` (`bootstrap`) sukses
