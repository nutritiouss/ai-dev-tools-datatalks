# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies for CGO (needed for SQLite driver, even though we use PostgreSQL)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (CGO enabled for SQLite driver compatibility)
RUN CGO_ENABLED=1 GOOS=linux go build -a -o server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
