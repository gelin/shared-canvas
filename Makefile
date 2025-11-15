.PHONY: build
build: build-go

.PHONY: build-go
build-go:
	go build ./cmd/server

.PHONY: run
run:
	./server
