# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY api/go.mod api/go.sum /app/
WORKDIR /app/api
RUN go mod download

# Copy the entire API project and build it
COPY api/ /app/api/
RUN go build -o server cmd/server/main.go

# Runtime stage
FROM ubuntu:latest

WORKDIR /root/

# Install necessary dependencies
RUN apt-get update && apt-get install -y ca-certificates

# Copy built server from builder stage
COPY --from=builder /app/api/server .
COPY api/cmd/server/.env /root/.env

# Expose the backend port
EXPOSE 3000

# Run the server
CMD ["./server"]
