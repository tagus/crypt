.PHONY: test

test:
	go test ./internal/...

crypt:
	go build -o crypt ./cmd/crypt/main.go
