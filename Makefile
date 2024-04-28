FORLINUX=GOOS=linux GOARCH=amd64

test:
	go test ./internal/...

crypt:
	go build -o crypt ./cmd/crypt/main.go

crypt.linux:
	$(FORLINUX) go build -o crypt.linux ./cmd/crypt/main.go

install: crypt
	mv crypt ${GOBIN}/crypt

clean:
	rm -f crypt

tidy:
	go mod tidy

.PHONY: clean crypt test install tidy
