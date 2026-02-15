FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o audit-service ./cmd/api

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/audit-service .

EXPOSE 8080

CMD ["./audit-service"]
