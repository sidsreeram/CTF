# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build it
COPY . . 
RUN go build -o server ./api/cmd/server/main.go
# Add this line in your Dockerfile
COPY template /app/template

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy built server from the builder stage
COPY --from=builder /app/server .
COPY api/cmd/server/.env .env

# Expose the backend port
EXPOSE 3000

# Run the server
CMD ["./server"]
