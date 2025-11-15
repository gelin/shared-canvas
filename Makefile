.PHONY: build
build: build-go

.PHONY: build-go
build-go:
	go build ./cmd/shared-canvas-server

.PHONY: run
run:
	./shared-canvas-server
