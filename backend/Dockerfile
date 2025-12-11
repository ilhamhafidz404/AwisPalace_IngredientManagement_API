# Stage 1 - Build binary
FROM golang:1.24 AS builder

WORKDIR /app

# Salin dependency file
COPY go.mod go.sum ./
RUN go env -w GOPROXY=direct && go mod download
# COPY .env ./

# Salin semua source code
COPY . .

# Build file binary
RUN go build -o main .

# Stage 2 - Minimal image
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
