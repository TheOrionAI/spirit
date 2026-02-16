.PHONY: build install test clean release

VERSION := 0.1.0
BINARY := spirit

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/$(BINARY) ./cmd/spirit

install:
	go install -ldflags "-X main.Version=$(VERSION)" ./cmd/spirit

test:
	go test -v ./...

clean:
	rm -rf bin/ dist/

release:
	goreleaser release --snapshot --skip-publish --rm-dist

# Quick dev build
dev:
	go run ./cmd/spirit

# Dependencies
deps:
	go mod tidy
	go mod download

# Lint
check:
	go vet ./...
	gofmt -w .
