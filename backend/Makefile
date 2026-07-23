AIR := $(shell command -v air 2>/dev/null || echo "$$HOME/go/bin/air")

.PHONY: dev dev-hot seed test tidy fmt build

dev:
	go run ./cmd/api

dev-hot:
	$(AIR) -c .air.toml

seed:
	go run ./cmd/seed

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w cmd internal

build:
	go build -o bin/api ./cmd/api
