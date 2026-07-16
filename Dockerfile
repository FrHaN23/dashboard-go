# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache build-base musl-dev

WORKDIR /app


COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
COPY openapi.yaml /app/openapi.yaml

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /app/backend-bin .

# ---------- Runtime stage ----------
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/backend-bin ./backend
COPY --from=builder /app/openapi.yaml ./openapi.yaml

RUN chmod +x ./backend

EXPOSE 5050

CMD ["./backend"]
