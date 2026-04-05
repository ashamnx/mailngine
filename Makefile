.PHONY: build run server worker migrate-up migrate-down sqlc test lint clean docker-build docker-up docker-down

# Build all binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/worker ./cmd/worker
	go build -o bin/migrate ./cmd/migrate

# Run the API server
run: server

server:
	go run ./cmd/server

# Run the background worker
worker:
	go run ./cmd/worker

# Database migrations
migrate-up:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir internal/db/migrations -seq $$name

# Generate sqlc code
sqlc:
	sqlc generate

# Run tests
test:
	go test -race -cover ./...

# Run linter
lint:
	golangci-lint run ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/

# Docker
docker-build:
	docker build -t mailngine .

docker-up:
	docker compose up -d

docker-down:
	docker compose down
