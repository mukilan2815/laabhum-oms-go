# Stage 1: Build the binary
FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o broker-adapter ./cmd/adapter

# Stage 2: Create the runtime container
FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/broker-adapter .

EXPOSE 8082

ENTRYPOINT ["/app/broker-adapter"]
