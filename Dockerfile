FROM golang:1.24-alpine

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

EXPOSE 4211

# Build the application
RUN go build -o main .

# Run the application
CMD ["./main"] 