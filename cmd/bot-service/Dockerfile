# syntax=docker/dockerfile:1

# Base Image
FROM golang:1.20

WORKDIR /app

# Install project dependencies
RUN apt-get update && apt-get install -y iputils-ping && apt install -y make
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy rest of source code files
COPY . ./
COPY ./.env ./.env

# Build project
RUN CGO_ENABLED=0 GOOS=linux go build -o build/bot/main ./cmd/bot-service/main.go 

# Other configs
VOLUME /app/logs

# Run the Discord bot service application
CMD sh -c "make run-aztebot-bot-service"