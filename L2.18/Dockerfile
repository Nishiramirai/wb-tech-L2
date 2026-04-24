# Stage 1: Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/calendar-app ./cmd/server

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /bin/calendar-app .

COPY --from=builder /app/config ./config

EXPOSE 8080
ENV CONFIG_PATH=config/local.yaml
CMD ["./calendar-app"]