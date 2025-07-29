# Development Dockerfile
FROM golang:1.24.2-alpine

# Install dependencies for development
RUN apk add --no-cache git ca-certificates tzdata

# Install tern migration tool
RUN go install github.com/jackc/tern/v2@latest

# Set timezone
ENV TZ=America/Sao_Paulo

WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 3333

# Default command for development (can be overridden in docker-compose)
CMD ["go", "run", "cmd/server/main.go"]
