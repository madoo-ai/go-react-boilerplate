SHELL := /bin/bash
.DEFAULT_GOAL := help
VERSION ?= 0.1.0-dev

ENV_RUN = set -a; [ -f .env ] && . ./.env; set +a;

.PHONY: help setup install env dev dev-api dev-web webui build test test-e2e test-e2e-install \
	lint vet typecheck i18n-check fmt fmt-check db-up db-down db-migrate db-seed db-reset \
	docker-up docker-down docker-logs skills skills-update clean

help:
	@printf '%s\n' \
	  'setup                Install deps, create .env, start Postgres/Mailpit, migrate and seed' \
	  'dev                  Go API on :8080 + Vite HMR on :5173' \
	  'build                Build the React app and the Go binaries into bin/' \
	  'test                 Go unit tests + Vitest' \
	  'test-e2e             Playwright smoke suite (needs the app running on :3000)' \
	  'lint                 go vet, gofmt check, TypeScript and i18n parity' \
	  'format               gofmt the Go trees' \
	  'db-up                Start Postgres and Mailpit' \
	  'db-down              Stop Postgres and Mailpit' \
	  'db-migrate           Apply embedded migrations' \
	  'db-seed              Seed roles, permissions and demo users' \
	  'db-reset             Delete the database volume and rebuild it (destructive)' \
	  'docker-up            Build and start the full Compose stack' \
	  'docker-down          Stop the Compose stack (keeps the data volume)' \
	  'skills               List the installed agent skills' \
	  'skills-update        Update the vendor agent skills' \
	  'clean                Remove bin/ and the built React assets'

install:
	cd web && npm ci --no-audit --no-fund

env:
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@if grep -q "replace-me-with-a-32-byte-secret" .env; then \
	  secret=$$(openssl rand -base64 32); \
	  sed -i "s|replace-me-with-a-32-byte-secret|$$secret|" .env; \
	  echo "Generated AUTH_SECRET in .env"; \
	fi

setup: install env db-up db-migrate db-seed webui
	@echo "Setup complete. Run: make dev (http://localhost:5173)"

dev: db-up
	@bash scripts/dev.sh

dev-api: db-up
	@$(ENV_RUN) go run ./cmd/server

dev-web:
	cd web && npm run dev

webui:
	cd web && npm run build

build: webui
	mkdir -p bin
	go build -ldflags "-X go-react-boilerplate/internal/xmodule.Version=$(VERSION)" -o bin/ ./cmd/...

test:
	go test ./cmd/... ./internal/...
	cd web && npm test

test-e2e:
	cd web && npx playwright test

test-e2e-install:
	cd web && npx playwright install chromium

lint: vet fmt-check typecheck i18n-check
	@if command -v golangci-lint >/dev/null 2>&1; then \
	  golangci-lint run ./cmd/... ./internal/...; \
	else \
	  echo "golangci-lint not installed; skipped (install: https://golangci-lint.run)"; \
	fi

vet:
	go vet ./cmd/... ./internal/...

typecheck:
	cd web && npm run typecheck

i18n-check:
	cd web && npm run check:i18n

fmt:
	gofmt -w cmd internal

fmt-check:
	@out=$$(gofmt -l cmd internal); if [ -n "$$out" ]; then echo "unformatted:"; echo "$$out"; exit 1; fi

db-up:
	@if pg_isready -h 127.0.0.1 -p 54330 -U postgres >/dev/null 2>&1; then \
	  echo "Using the profile PostgreSQL on 127.0.0.1:54330"; \
	else \
	  docker compose up -d --wait db mailpit; \
	fi

db-down:
	docker compose down

db-migrate:
	@$(ENV_RUN) go run ./cmd/migrate

db-seed:
	@$(ENV_RUN) go run ./cmd/seed

db-reset:
	docker compose down -v
	$(MAKE) db-up db-migrate db-seed

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

skills:
	npx --yes skills list

skills-update:
	npx --yes skills update

clean:
	rm -rf bin web/dist-ssr web/test-results web/playwright-report
