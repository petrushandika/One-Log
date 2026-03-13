# Variables
BACKEND_DIR=Backend
FRONTEND_DIR=Frontend

.PHONY: all setup generate migrate build dev docker-up docker-down

all: setup generate build

# 🛠️ Setup Project
setup:
	@echo "Installing dependencies..."
	cd $(BACKEND_DIR) && go mod download
	cd $(FRONTEND_DIR) && npm install --legacy-peer-deps

# 🏗️ Code Generation (Mocks, Protobuf, etc.)
generate:
	@echo "Generating code..."
	cd $(BACKEND_DIR) && go generate ./...

# 🗄️ Database Migration
migrate:
	@echo "Running database migrations..."
	# This will trigger the AutoMigrate in your Go code
	cd $(BACKEND_DIR) && go run cmd/api/main.go --migrate

# 🔨 Build All
build:
	@echo "Building applications..."
	cd $(BACKEND_DIR) && go build -o bin/api ./cmd/api
	cd $(FRONTEND_DIR) && npm run build

# 🚀 Run Development
dev:
	@echo "Starting development servers..."
	# Run backend in background
	cd $(BACKEND_DIR) && go run cmd/api/main.go &
	# Run frontend
	cd $(FRONTEND_DIR) && npm run dev

# 🐳 Docker Operations
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

# 🧹 Clean Up
clean:
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(FRONTEND_DIR)/dist
