# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o server .

# Expose the port the server will use
EXPOSE 8082

# Run the application
CMD ["./server"]
