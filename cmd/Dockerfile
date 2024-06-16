# syntax=docker/dockerfile:1

# Base Image
FROM golang:latest

WORKDIR /app

# Install project dependencies
RUN apt-get update && apt-get install -y iputils-ping
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of source code files
COPY . ./
COPY ./.env ./.env

# Build project
RUN CGO_ENABLED=0 GOOS=linux go build -o build/bot/main ./cmd/main.go 

# Run the Discord bot service application
CMD sh -c "./build/bot/main"