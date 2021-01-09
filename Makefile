.PHONY: test

test:
	go test ./internal/...

crypt: 
	go build -o crypt ./cmd/crypt/main.go

install: crypt
	mv crypt ${GOBIN}/crypt

clean:
	rm -f crypt

.PHONY: clean crypt test install