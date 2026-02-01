# build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o shortener cmd/server/main.go

# final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/shortener .
EXPOSE 8080
ENTRYPOINT ["./shortener"]
