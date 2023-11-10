migrate-up:
	echo "Running migrations..." && goose -dir ./internal/bot-service/data/migrations up

run-services:
	./build/main