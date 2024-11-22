# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/employee/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and necessary files
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/db/migrations ./db/migrations

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]