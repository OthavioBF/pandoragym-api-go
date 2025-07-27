FROM golang:1.24.2-alpine AS base

# Install dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# For development - we'll run commands directly
# The actual command will be specified in docker-compose.yml

# Expose port
EXPOSE 3333

# Default command (can be overridden in docker-compose)
CMD ["go", "run", "cmd/server/main.go"]
