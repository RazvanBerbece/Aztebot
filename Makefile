# LOCAL DEVELOPMENT UTILITY SHELL APPS
migrate-up:
	sql-migrate up -config=local.dbconfig.yml -env="staging"

migrate-up-dry:
	sql-migrate up -config=local.dbconfig.yml -env="staging" -dryrun

run-all:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-env:
	openssl base64 -A -in .prod.env -out .env.out


# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main

run-azteradio-service:
	./build/azteradio/main