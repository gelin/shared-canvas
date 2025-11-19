.PHONY: build
build: build-webapp build-server

.PHONY: install-server
install-server:
	go mod download

.PHONY: build-server
build-server: install-server
	go build ./cmd/shared-canvas-server

.PHONY: install-webapp
install-webapp:
	cd webapp && npm ci || npm install

.PHONY: build-webapp
build-webapp: install-webapp
	cd webapp && npm run build

.PHONY: clean
clean: clean-webapp clean-server

.PHONY: clean-webapp
clean-webapp:
	rm -rf ./cmd/shared-canvas-server/web-dist/*

.PHONY: clean-server
clean-server:
	rm -f ./shared-canvas-server

.PHONY: run
run:
	./shared-canvas-server

.PHONY: dev
# Run Go server (port 8080) and Vite dev server (port 5173) in two terminals:
# 1) make run
# 2) cd webapp && npm run dev
# The Vite dev server proxies /api to http://localhost:8080
dev:
	@echo "Open two terminals: 'make run' and 'cd webapp && npm run dev'"
