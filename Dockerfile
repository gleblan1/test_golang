# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy swagger documentation
COPY docs/swagger.json /app/docs/swagger.json

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Build Worker binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# API stage
FROM alpine:latest AS api

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./api"]

# Worker stage
FROM alpine:latest AS worker

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/worker .
COPY --from=builder /app/configs ./configs

CMD ["./worker"] 