FROM golang:1.21 AS builder

# Install build dependencies for SQLite
RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev git && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o godoor .

FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/godoor .

# Create directory for database
RUN mkdir -p /app/data

# Set working directory
WORKDIR /app

# Copy binary to working directory
COPY --from=builder /app/godoor .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./godoor"]