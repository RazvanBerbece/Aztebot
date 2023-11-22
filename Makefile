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
	openssl base64 -A -in .prod.env -out .env.out

update-music-conf:
	openssl base64 -A -in cmd/azteradio-service-rawon/.config.prod.env -out cmd/azteradio-service-rawon/.base64.prod.env
	openssl base64 -A -in cmd/aztemusic-service-rawon/1/.config.prod.env -out cmd/aztemusic-service-rawon/1/.base64.prod.env


# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main

run-aztemusic-service:
	./build/aztemusic/main

run-azteradio-orchestrator-service:
	./build/azteradio-orchestrator-service/main