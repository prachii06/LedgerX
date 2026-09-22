# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application as a statically linked Linux executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ledgerx-server ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/ledgerx-server .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./ledgerx-server"]
