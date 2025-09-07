# 1. Base Go image
FROM golang:1.22-alpine

# 2. Set working directory
WORKDIR /app

# 3. Install git + certificates
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# 4. Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./

# 5. Download deps
RUN go mod download

# 6. Copy project files
COPY . .

# 7. Build Go binary
RUN go build -o main .

# 8. Copy SSL cert
COPY isrgrootx1.pem /app/isrgrootx1.pem

# 9. Expose port (must match your Go server port)
EXPOSE 6001

# 10. Run the app
CMD ["./main"]
