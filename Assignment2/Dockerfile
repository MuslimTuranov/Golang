FROM golang:1.25.6-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/myapp ./cmd/api

FROM alpine:3.22

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /out/myapp /app/myapp
COPY --from=builder /src/database/migrations /app/database/migrations

EXPOSE 8080

USER appuser

ENTRYPOINT ["/app/myapp"]
