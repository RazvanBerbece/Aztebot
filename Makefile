# LOCAL DEVELOPMENT UTILITY SHELL APPS
migrate-up:
	sql-migrate up -config=local.dbconfig.yml -env="staging"

migrate-up-dry:
	sql-migrate up -config=local.dbconfig.yml -env="staging" -dryrun

up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-conf-jar:
	openssl base64 -A -in ./cmd/aztemusic-service-jar/aztemusic-service-1/config.prod.txt -out ./cmd/aztemusic-service-jar/aztemusic-service-1/base64.config.prod.txt
	openssl base64 -A -in ./cmd/aztemusic-service-jar/aztemusic-service-2/config.prod.txt -out ./cmd/aztemusic-service-jar/aztemusic-service-2/base64.config.prod.txt
	openssl base64 -A -in ./cmd/aztemusic-service-jar/aztemusic-service-3/config.prod.txt -out ./cmd/aztemusic-service-jar/aztemusic-service-3/base64.config.prod.txt
	openssl base64 -A -in ./cmd/aztemusic-service-jar/azteradio-service/config.prod.txt -out ./cmd/aztemusic-service-jar/azteradio-service/base64.config.prod.txt

update-env:
	openssl base64 -A -in .prod.env -out base64.prod.env.out
	openssl base64 -A -in cmd/aztemusic-service/1/.prod.env -out cmd/aztemusic-service/1/base64.prod.env.out

# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main

run-aztemusic-bot-service:
	./build/music/main