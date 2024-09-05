# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the backend source code to the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose the port on which your backend listens (adjust if needed)
EXPOSE 7000

# Run the backend application
CMD ["./main"]