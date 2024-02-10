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
RUN CGO_ENABLED=0 GOOS=linux go build -o build/bot/main ./cmd/bot-service/main.go 

# Other configs
VOLUME /app/logs

# Run the Discord bot service application
CMD ["go", "test", "-v", "./..."]