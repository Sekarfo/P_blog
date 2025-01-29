# Use the official Golang image
FROM golang:1.20 as builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy Go files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Use minimal image for final build
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY web/templates /root/templates
COPY web/static /root/static

EXPOSE 8080

CMD ["./main"]
