# Stage 1: Build the application
FROM golang:1.20 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o adapter ./cmd/adapter

# Stage 2: Create the final image
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/adapter .

EXPOSE 8082
CMD ["./adapter"]
