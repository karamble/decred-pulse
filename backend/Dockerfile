# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Copy source code (needed for go mod tidy to work)
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o decred-pulse-backend .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/decred-pulse-backend .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./decred-pulse-backend"]

