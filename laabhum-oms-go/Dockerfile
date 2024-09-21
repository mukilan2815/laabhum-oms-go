# Stage 1: Build the application
FROM golang:1.20 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o oms ./cmd/server

# Stage 2: Create the final image
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/oms .

EXPOSE 8081
CMD ["./oms"]
