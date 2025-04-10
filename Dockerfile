FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy Go modules files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Copy the .env file explicitly
# COPY ../.env .env

# Build the application
RUN go build -o wallet-service main.go

# Use a minimal final image
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/wallet-service .
# COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./wallet-service"]
