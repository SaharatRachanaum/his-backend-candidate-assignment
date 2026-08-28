# ==========================================
# Build Stage
# ==========================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server ./cmd/server

# ==========================================
# Final Production Stage
# ==========================================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Bangkok

# Run as non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose Go application port
EXPOSE 8080

CMD ["./server"]
