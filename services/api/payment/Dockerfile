# Use the official Golang image from Docker Hub with the correct version tag
FROM golang:1.22.5

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (if present) and download dependencies
COPY . .

RUN go get -d -v ./...

RUN go build -o main .

# Command to run the executable
CMD ["./main"]