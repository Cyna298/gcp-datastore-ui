# Default package manager for the frontend
PKG_MANAGER ?= pnpm

# Default port settings
# CAUTION: The server port is also used in the frontend ajaxPromise.ts file, make sure to change it there as well if you change it here and rebuild the frontend. This should be a smooth process in the future.
SERVER_PORT ?= 8080
EMU_PORT ?= 8081
PROJECT_ID ?= exactbuyer-api-staging

# Paths
BACKEND_PATH := backend
FRONTEND_PATH := webapp

RED := \033[0;31m
GREEN := \033[0;32m
NO_COLOR := \033[m

.PHONY: help setup-backend setup-frontend run-backend run-frontend build-backend build-frontend check-pkg-manager check-and-copy-out clean-run

run-backend:
	@echo "Starting the backend server on port $(PORT)..."
	@templ generate
	@tailwindcss -o public/styles.css
	@air --build.cmd "go build -o bin/service" --build.bin "bin/service --port $(SERVER_PORT)\
		--project $(PROJECT_ID)\
		--emuHost localhost:$(EMU_PORT)\
		--emuHostPath localhost:$(EMU_PORT)/datastore\
		--dsHost http://localhost:$(EMU_PORT)"

templ:
	templ generate --watch --proxy="http://localhost:8080"
tailwind:
	tailwindcss -o public/styles.css --watch
