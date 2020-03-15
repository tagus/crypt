.PHONY: test

test:
	go test ./internal/...

build:
	go build -o crypt ./cmd/crypt/main.go