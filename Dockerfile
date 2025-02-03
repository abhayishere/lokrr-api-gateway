FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy all project files into the container
COPY . .

# Download Go module dependencies
RUN go mod tidy

# Build the binary from the main.go file in the cmd folder
RUN go build -o api-gateway ./cmd/main.go

# Expose port 8080 for the API Gateway
EXPOSE 8080

# Command to run the built binary
CMD ["./api-gateway"]
