# LOCAL DEVELOPMENT UTILITY SHELL APPS
migrate-up:
	sql-migrate up -config=local.dbconfig.yml -env="staging"

migrate-up-dry:
	sql-migrate up -config=local.dbconfig.yml -env="staging" -dryrun

up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-env:
	openssl base64 -A -in .prod.env -out .base64.prod.env.out

update-music-conf:
	openssl base64 -A -in cmd/azteradio-service/config.prod.yml -out cmd/azteradio-service/base64.prod.yml.out
	openssl base64 -A -in cmd/aztemusic-service/1/config.prod.yml -out cmd/aztemusic-service/1/base64.prod.yml.out

# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main