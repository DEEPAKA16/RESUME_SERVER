# Use the correct Go version (>= 1.24.2)
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o app .

# Final lightweight image
FROM alpine:latest
WORKDIR /root/

# Copy compiled binary
COPY --from=builder /app/app .

# Run app
CMD ["./app"]
