# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire current directory (including the Go code) to the container's /app directory
COPY . /app

# Install the required packages for your Go application (if needed)
# RUN go get -d -v ./...

# Build your Go application
RUN go build -o main .

# Expose the port on which your Go application listens
EXPOSE 9000

# Run your Go application when the container starts
CMD ["./main"]
