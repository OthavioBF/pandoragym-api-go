# Development Dockerfile with Air hot reload
FROM golang:1.24.2-alpine

# Install dependencies for development
RUN apk add --no-cache git ca-certificates tzdata curl

# Install development tools
RUN go install github.com/jackc/tern/v2@latest
RUN go install github.com/air-verse/air@latest

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

# Default command for development with Air hot reload
CMD ["air", "-c", ".air.toml"]
