# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Stage 2: Create a minimal image
FROM alpine:latest  

# Add root certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary and the migrations folder
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

# Expose port
EXPOSE 8081

# Run the binary
CMD ["./main"]
