# Default package manager for the frontend
PKG_MANAGER ?= pnpm

# Default port settings
# CAUTION: The server port is also used in the frontend ajaxPromise.ts file, make sure to change it there as well if you change it here and rebuild the frontend. This should be a smooth process in the future.
SERVER_PORT ?= 8080
EMU_PORT ?= 8081
PROJECT_ID ?= project-id

# Paths
BACKEND_PATH := backend
FRONTEND_PATH := webapp

RED := \033[0;31m
GREEN := \033[0;32m
NO_COLOR := \033[m

.PHONY: help setup-backend setup-frontend run-backend run-frontend build-backend build-frontend check-pkg-manager check-and-copy-out clean-run

help:
	@echo "Makefile commands:"
	@echo "${GREEN}setup-backend${NO_COLOR}  - Install Go dependencies for the backend."
	@echo "${GREEN}setup-frontend${NO_COLOR} - Install frontend dependencies using $(PKG_MANAGER)."
	@echo "${GREEN}setup${NO_COLOR}          - Run setup-backend, setup-frontend and check-and-copy-out."
	@echo "${GREEN}run-backend${NO_COLOR}    - Start the backend server after checking and copying 'out' directory."
	@echo "${GREEN}run-frontend${NO_COLOR}   - Start the frontend server."
	@echo "${GREEN}build-backend${NO_COLOR}  - Build the backend application."
	@echo "${GREEN}build-frontend${NO_COLOR} - Build the frontend application."
	@echo "${GREEN}clean-run${NO_COLOR}      - Clean 'out' directories and perform a fresh frontend build."
	@echo "Use ${RED}PKG_MANAGER${NO_COLOR}=[pnpm|npm|yarn] to specify the frontend package manager."

setup-backend:
	@echo "Installing Go dependencies for the backend..."
	cd $(BACKEND_PATH) && go get ./...

setup-frontend: check-pkg-manager
	@echo "Installing frontend dependencies using $(PKG_MANAGER)..."
	cd $(FRONTEND_PATH) && $(PKG_MANAGER) install

check-pkg-manager:
	@which $(PKG_MANAGER) > /dev/null || (echo "$(PKG_MANAGER) is not installed. Please install it or specify another package manager using PKG_MANAGER=[pnpm|npm|yarn]" && exit 1)

check-and-copy-out:
	@if [ ! -d "$(BACKEND_PATH)/out" ]; then \
		echo "Checking for out directory in frontend..."; \
		if [ -d "$(FRONTEND_PATH)/out" ]; then \
			echo "Copying out directory from frontend to backend..."; \
			cp -r $(FRONTEND_PATH)/out $(BACKEND_PATH)/; \
		else \
			echo "out directory not found in frontend, building frontend..."; \
			$(MAKE) build-frontend; \
			cp -r $(FRONTEND_PATH)/out $(BACKEND_PATH)/; \
		fi \
	fi

run-backend: check-and-copy-out
run-backend:
	@echo "Starting the backend server on port $(PORT)..."
	cd $(BACKEND_PATH) && go run . -port $(SERVER_PORT)\
		-project $(PROJECT_ID)\
		-emuHost localhost:$(EMU_PORT)\
		-emuHostPath localhost:$(EMU_PORT)/datastore\
		-dsHost http://localhost:$(EMU_PORT)

run-frontend:
	@echo "Starting the frontend server..."
	cd $(FRONTEND_PATH) && $(PKG_MANAGER) start

build-backend:
	@echo "Building the backend application..."
	cd $(BACKEND_PATH) && go build -o bin/service .

build-frontend:
	@echo "Building the frontend application..."
	cd $(FRONTEND_PATH) && $(PKG_MANAGER) run build

clean-run:
	@echo "Cleaning out directories..."
	@if [ -d "$(FRONTEND_PATH)/out" ]; then rm -rf $(FRONTEND_PATH)/out; fi
	@if [ -d "$(BACKEND_PATH)/out" ]; then rm -rf $(BACKEND_PATH)/out; fi
	@echo "Building frontend..."
	@$(MAKE) build-frontend
	@echo "Copying out directory from frontend to backend..."
	@cp -r $(FRONTEND_PATH)/out $(BACKEND_PATH)/

setup: setup-backend setup-frontend check-and-copy-out
