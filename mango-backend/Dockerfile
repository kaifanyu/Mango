# Stage 1: Build the application
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency resolution
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
