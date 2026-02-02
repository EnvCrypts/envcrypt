FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd


# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/goose /bin/goose
COPY --from=builder /app/server .

COPY --from=builder /app/internal/db/migrations ./migrations
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh && touch .env

EXPOSE 8080

CMD ["./entrypoint.sh"]
