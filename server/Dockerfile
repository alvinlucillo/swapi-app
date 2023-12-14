# Create dockerfile for the server with main.go as the entrypoint
# Use a base image
FROM golang:alpine

# Set the working directory
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o server ./cmd/main.go

# Set the entrypoint for the container
ENTRYPOINT ["./server"]
