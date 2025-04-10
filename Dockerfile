FROM golang:1.23-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git

# Create and set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with version information
ARG VERSION=dev
ARG BUILD_DATE
ARG COMMIT

RUN go build -ldflags="-X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -X 'main.Commit=${COMMIT}'" -o /go/bin/radar ./cmd/radar

# Create a smaller final image
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/radar /usr/local/bin/radar

# Create a directory for data
RUN mkdir -p /app/data

# Set working directory
WORKDIR /app

# Create entrypoint
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
