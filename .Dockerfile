# Use the latest official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080
EXPOSE 8082

# Command to run the executable
CMD ["./main"]