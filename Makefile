.PHONY: help start stop restart logs logs-dcrd logs-backend logs-frontend build clean status shell-dcrd shell-backend backup restore

# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Default RPC credentials (override with env vars)
DCRD_RPC_USER ?= decred
DCRD_RPC_PASS ?= decredpass

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

start: ## Start all services (dcrd + backend + frontend)
	@echo "Starting Decred Dashboard stack..."
	docker compose up -d
	@echo "Services started! Dashboard will be available at http://localhost:3000"
	@echo "Note: dcrd needs to sync the blockchain on first run (may take several hours)"

stop: ## Stop all services
	@echo "Stopping all services..."
	docker compose down

restart: ## Restart all services
	@echo "Restarting all services..."
	docker compose restart

logs: ## View logs from all services
	docker compose logs -f

logs-dcrd: ## View dcrd logs
	docker compose logs -f dcrd

logs-backend: ## View backend logs
	docker compose logs -f backend

logs-frontend: ## View frontend logs
	docker compose logs -f frontend

build: ## Build/rebuild all images
	@echo "Building images..."
	docker compose build --no-cache

status: ## Show status of all services
	@echo "=== Service Status ==="
	docker compose ps
	@echo ""
	@echo "=== dcrd Info ==="
	@docker exec decred-pulse-dcrd dcrctl \
		--rpcuser=$(DCRD_RPC_USER) \
		--rpcpass=$(DCRD_RPC_PASS) \
		--rpcserver=127.0.0.1:9109 \
		--rpccert=/certs/rpc.cert \
		getinfo 2>/dev/null || echo "dcrd not ready yet..."

clean: ## Stop and remove all containers, networks, and volumes (WARNING: deletes blockchain data!)
	@echo "WARNING: This will delete the blockchain data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose down -v; \
		echo "Cleaned up!"; \
	fi

shell-dcrd: ## Open shell in dcrd container
	docker exec -it decred-pulse-dcrd /bin/sh

shell-backend: ## Open shell in backend container
	docker exec -it decred-pulse-backend /bin/sh

dcrctl: ## Run dcrctl command (usage: make dcrctl CMD="getblockcount")
	@docker exec decred-pulse-dcrd dcrctl \
		--rpcuser=$(DCRD_RPC_USER) \
		--rpcpass=$(DCRD_RPC_PASS) \
		--rpcserver=127.0.0.1:9109 \
		--rpccert=/certs/rpc.cert \
		$(CMD)

sync-status: ## Check blockchain sync status
	@echo "=== Blockchain Sync Status ==="
	@docker exec decred-pulse-dcrd dcrctl \
		--rpcuser=$(DCRD_RPC_USER) \
		--rpcpass=$(DCRD_RPC_PASS) \
		--rpcserver=127.0.0.1:9109 \
		--rpccert=/certs/rpc.cert \
		getblockchaininfo 2>/dev/null || echo "dcrd not ready yet..."

peers: ## Show connected peers
	@echo "=== Connected Peers ==="
	@docker exec decred-pulse-dcrd dcrctl \
		--rpcuser=$(DCRD_RPC_USER) \
		--rpcpass=$(DCRD_RPC_PASS) \
		--rpcserver=127.0.0.1:9109 \
		--rpccert=/certs/rpc.cert \
		getpeerinfo 2>/dev/null | grep -E '"addr"|"id"' || echo "dcrd not ready yet..."

backup: ## Backup blockchain data
	@echo "Creating backup of dcrd data..."
	@mkdir -p backups
	docker run --rm \
		-v umbrel-decred-dash-main_dcrd-data:/data \
		-v $(PWD)/backups:/backup \
		alpine tar czf /backup/dcrd-backup-$$(date +%Y%m%d-%H%M%S).tar.gz -C /data .
	@echo "Backup created in backups/"

restore: ## Restore blockchain data from backup (usage: make restore BACKUP=backups/dcrd-backup-xxx.tar.gz)
	@if [ -z "$(BACKUP)" ]; then \
		echo "Usage: make restore BACKUP=backups/dcrd-backup-xxx.tar.gz"; \
		exit 1; \
	fi
	@echo "Restoring from $(BACKUP)..."
	@echo "WARNING: This will overwrite existing blockchain data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose down; \
		docker run --rm \
			-v umbrel-decred-dash-main_dcrd-data:/data \
			-v $(PWD):/backup \
			alpine sh -c "rm -rf /data/* && tar xzf /backup/$(BACKUP) -C /data"; \
		docker compose up -d; \
		echo "Restored!"; \
	fi

setup: ## Initial setup (copy env.example to .env)
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Created .env file. Please edit it with your RPC credentials!"; \
		echo "Then run: make start"; \
	else \
		echo ".env file already exists"; \
	fi

update: ## Update all services to latest
	@echo "Rebuilding all services..."
	docker compose build --no-cache
	@echo "Updated! Run 'make restart' to apply changes"

update-dcrd: ## Update dcrd to latest from source
	@echo "Rebuilding dcrd from latest source..."
	docker compose build --no-cache dcrd
	docker compose up -d dcrd
	@echo "dcrd updated!"

build-dcrd: ## Build dcrd (usage: make build-dcrd VERSION=release-v2.0.6)
	@if [ -n "$(VERSION)" ]; then \
		echo "Building dcrd $(VERSION)..."; \
		DCRD_VERSION=$(VERSION) docker compose build --no-cache dcrd; \
	else \
		echo "Building dcrd from master..."; \
		docker compose build dcrd; \
	fi

dev-backend: ## Run backend in development mode (outside Docker)
	cd backend && go run main.go

dev-frontend: ## Run frontend in development mode (outside Docker)
	cd frontend && npm run dev

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

install-backend: ## Install backend dependencies
	cd backend && go mod download

