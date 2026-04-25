.PHONY: test run build migrate-up migrate-down

# Run all unit tests
test:
	go test -v ./...

# Run the backend server
run:
	go run cmd/main.go

# Build the binary
build:
	go build -o bin/server cmd/main.go

# Run migrations (if you have the CLI installed)
migrate-up:
	migrate -path db/migrations -database "postgres://postgres:root@localhost:8080/postgres?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://postgres:root@localhost:8080/postgres?sslmode=disable" down
