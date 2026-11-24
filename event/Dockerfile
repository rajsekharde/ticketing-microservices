# Build Stage

# Use a lightweight alpine-based Go image. Stage named as builder
FROM golang:1.25-alpine AS builder


# Set environment variables for the Go build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# CGO_ENABLED=0 : Produces a pure Go binary that's not dependent on C libraries
# GOOS=linux : target OS
# GOARCH=amd64 : target architecture


# Set working directory inside the container to /build
WORKDIR /build

# Creates /build directory if it doesn't exist
# Sets /build as current working directory


# Install ca-certificates (needed for outbound HTTPS/gRPC TLS if used)
RUN apk add --no-cache ca-certificates


# Copy only go.mod, go.sum first 
COPY go.mod go.sum ./


# Download all dependencies
RUN go mod download


# Copy the rest of the source code into the container
COPY . .


# Build the Go app and output the binary at /app
RUN go build -o /event ./event


# Run Stage

# Start a fresh, minimal alpine image. This is the final runtime image
FROM alpine:3.21 AS final


RUN apk add --no-cache ca-certificates

# Copy compiled binaries from builder stage into this stage
COPY --from=builder /event /bin/event


# Documents that app listens on port 8000
EXPOSE 50052


# Default command when container starts. Runs the Go binary
CMD ["bin/event"]