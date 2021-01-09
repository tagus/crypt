.PHONY: test

test:
	go test ./internal/...

crypt: cmd/crypt/main.go
	go build -o crypt ./cmd/crypt/main.go

clean:
	rm -f crypt
