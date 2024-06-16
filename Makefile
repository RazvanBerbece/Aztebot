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

update-jar-conf:
	openssl base64 -A -in cmd/azteradio-service-java/.prod.config.txt -out .prod.jar-config.txt


# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main

run-azteradio-service:
	./build/azteradio/main