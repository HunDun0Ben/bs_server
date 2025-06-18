# Start from the official golang image
FROM golang:1.20-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller base image for the final image
FROM alpine:latest

# Copy the binary from builder
COPY --from=0 /app/main /app/main

# Set working directory
WORKDIR /app

# Expose port (adjust if needed)
EXPOSE 8080

# Run the binary
CMD ["./main"]