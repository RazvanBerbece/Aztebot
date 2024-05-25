# syntax=docker/dockerfile:1

# Base Image
FROM golang:latest

WORKDIR /app

# Install project dependencies
RUN apt-get update
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of source code files
COPY . ./
COPY ./.env ./.env

# Build project
RUN CGO_ENABLED=0 GOOS=linux go build -o build/bot/main ./cmd/main.go 

# No entrypoint here as the entrypoint is defined in the docker-compose file for CI.
# This is because before the tests are run, some dependency services need to be booted up.