# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build for the target architecture (ARM64 for Orange Pi 5)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/bot ./cmd/bot/main.go

# Production stage
FROM alpine:3.18

# We need chroot which is in coreutils or busybox (included in alpine)
RUN apk add --no-cache coreutils

WORKDIR /app
COPY --from=builder /app/bot .

# ENV defaults
ENV HOST_ROOT=/host

ENTRYPOINT ["./bot"]
