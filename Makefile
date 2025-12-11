.PHONY: dev build run test lint clean

dev:
	go mod tidy
	go run ./cmd/auth

build:
	@mkdir -p tmp/bin
	go build -o tmp/bin/auth ./cmd/auth

run:
	go run ./cmd/auth

lint:
	golangci-lint run

clean:
	rm -rf tmp/

test:
	go test ./... -v

