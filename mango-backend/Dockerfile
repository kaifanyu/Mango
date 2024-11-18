# Use an official Go image
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum for dependency resolution
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Expose the application port
EXPOSE 8080

# Use Air to run the application
CMD ["air", "-c", ".air.toml"]
