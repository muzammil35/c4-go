# Use the official Golang image as build stage
FROM golang:1.23-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files and templates
COPY --from=builder /app/templates ./templates

# Expose port 8080
EXPOSE 8080

# Set environment variable for port
ENV PORT=8080

# Run the application
CMD ["./main"]