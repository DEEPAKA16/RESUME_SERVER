# 1. Use official lightweight Go image
FROM golang:1.22-alpine

# 2. Set working directory inside container
WORKDIR /app

# 3. Install required tools (git + CA certificates for SSL/TiDB)
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# 4. Copy dependency files first (for caching)
COPY go.mod go.sum ./

# 5. Download dependencies
RUN go mod download

# 6. Copy all source code into container
COPY . .

# 7. Build the Go binary
RUN go build -o main .

# 8. Copy SSL certificate (make sure isrgrootx1.pem is in your repo root)
COPY isrgrootx1.pem /app/isrgrootx1.pem

# 9. Expose your app’s port (change if your server runs on another port)
EXPOSE 6001

# 10. Run the built binary
CMD ["./main"]
