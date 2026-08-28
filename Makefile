.PHONY: all build run test clean docker-up docker-down docker-logs

# Application binary name
APP_NAME=hospital-middleware-server

all: test build

build:
	@echo "Building Go binary..."
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	@echo "Running server locally..."
	go run ./cmd/server

test:
	@echo "Running all unit tests..."
	go test -v ./...

docker-up:
	@echo "Starting all services with Docker Compose..."
	docker compose up --build -d

docker-down:
	@echo "Stopping all services..."
	docker compose down

docker-logs:
	@echo "Viewing logs..."
	docker compose logs -f

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	docker compose down -v
