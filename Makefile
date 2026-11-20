.DEFAULT_GOAL := build
.PHONY: build test test-race vet lint fmt watch clean

build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

watch:
	go tool air

clean:
	rm -rf bin tmp
