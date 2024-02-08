# LOCAL DEVELOPMENT UTILITY SHELL APPS
migrate-up:
	sql-migrate up -config=local.dbconfig.yml -env="staging"

migrate-up-dry:
	sql-migrate up -config=local.dbconfig.yml -env="staging" -dryrun

up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

ci:
	docker-compose -f docker-compose.ci.yml up --remove-orphans --build --exit-code-from integration-test-bot-service

update-env:
	openssl base64 -A -in .prod.env -out base64.prod.env.out

update-ci-env:
	openssl base64 -A -in .env -out base64.ci.env.out

# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main