# Build stage
FROM golang:1.24-alpine AS builder

# Install git and make for fetching dependencies and building
RUN apk add --no-cache git make bash

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
ARG VERSION
RUN make build version=${VERSION}

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/dist/go-project-template ./main
COPY --from=builder /app/config ./config

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
