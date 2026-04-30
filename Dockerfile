# Use official Go image
FROM golang:1.22-alpine

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all project files
COPY . .

# Build app
RUN go build -o main .

# Expose port
EXPOSE 8080

# Run app
CMD ["./main"]