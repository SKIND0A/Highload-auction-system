FROM golang:1.23.0-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/auction ./cmd/main.go

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 10001 appuser

COPY --chown=appuser:appuser --from=builder /app/bin/auction /app/auction
COPY --chown=appuser:appuser --from=builder /app/migrations /app/migrations
COPY --chown=appuser:appuser --from=builder /app/web /app/web

USER appuser
EXPOSE 8080
CMD ["/app/auction"]