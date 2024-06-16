migrate-up-cli:
	go run ./internal/bot-service/data/migrations/main.go

migrate-up:
	echo "Running migrations..." && goose -dir ./internal/bot-service/data/migrations up

run-services:
	./build/main