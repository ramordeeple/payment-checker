BINARY=payment-checker

.PHONY: build run test

build:
	go build -o $(BINARY) ./cmd/payment-checker/main.go

run:
	go run ./cmd/payment-checker/main.go

test:
	go test ./internal/... -v # Удобно на случай новых тестов
