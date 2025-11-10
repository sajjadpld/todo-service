.PHONY: env
env:
	@cp .env.example .env
	@echo ".env generated"

.PHONY: swag
swag:
	@echo "swagger generating..."
	@swag init -g ./cmd/main.go -q -o ./docs
	@echo "swagger generated"

.PHONY: tests
tests:
	@go test ./internal/adapter/repository -run TestTodoRepository_Create -v
	@go test ./internal/core/usecase -run TestTodoUsecase_Create -v
	@echo "TESTS WERE DONE"
