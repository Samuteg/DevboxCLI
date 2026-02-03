# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o devboxcli .

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/devboxcli .
# Garante que o binário seja executável
RUN chmod +x devboxcli
ENTRYPOINT ["./devboxcli"]