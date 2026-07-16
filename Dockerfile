# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache build-base musl-dev

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /app/backend-bin .

# ---------- Runtime stage ----------
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/backend-bin ./backend

RUN chmod +x ./backend

EXPOSE 5050

CMD ["./backend"]
