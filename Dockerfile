FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Build the Go application
RUN go build -o scheme_parser main.go compiler.go scheme_parser.go

# Expose the application port
EXPOSE 1234

# Run the application with the default port
CMD ["./scheme_parser"]
