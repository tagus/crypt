FORLINUX=GOOS=linux GOARCH=amd64
GO=CGO_ENABLED=1 go

test:
	$(GO) test ./internal/...

crypt:
	$(GO) build -o crypt ./cmd/crypt/main.go

crypt.linux:
	$(FORLINUX) $(GO) build -o crypt.linux ./cmd/crypt/main.go

install:
	$(GO) install ./cmd/crypt

clean:
	rm -f crypt

tidy:
	$(GO) mod tidy

.PHONY: clean crypt test install tidy
