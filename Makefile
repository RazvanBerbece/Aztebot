# LOCAL DEVELOPMENT UTILITY SHELL APPS
migrate-up:
	sql-migrate up -config=local.dbconfig.yml -env="staging"

migrate-up-dry:
	sql-migrate up -config=local.dbconfig.yml -env="staging" -dryrun

migrate-rollback:
	sql-migrate down -config=local.dbconfig.yml -env="staging"

up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

ci:
	docker compose down -v
	docker-compose -f docker-compose.ci.yml up --remove-orphans --force-recreate --build --exit-code-from integration-test-bot-service

update-envs:
	openssl base64 -A -in .prod.env -out base64.prod.env.out
	openssl base64 -A -in .env -out base64.ci.env.out

# APP STARTUP SHELL APPS
run-aztebot-bot-service:
	./build/bot/main