# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build frontend
RUN apk add --no-cache nodejs npm && \
    cd frontend && npm install && npm run build

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notifapi .

# Final stage
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates wget tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary and frontend build
COPY --from=builder /app/notifapi .
COPY --from=builder /app/frontend/dist ./frontend/dist

# Create logs directory
RUN mkdir -p logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 10887

CMD ["./notifapi"]