# Aztebot
The ambitious and robust Discord bot which powers the OTA (Ordinul Templierilor Azteci) Discord community. Written in Go.

Composing services:
- `mysql-db` (Containerised MySQL service for the data layer)
- `bot-service` (Handles Discord interactions like new messages, slash commands, join events, reaction adding or removing, etc.)

# How to Run
## Prerequisites
In order to run the application, a few prerequisites must be met.
1. Have the repository cloned locally.
2. Have Docker installed.
3. Have a fully-configured `.env` file saved in the root of the repository. (contact [@RazvanBerbece](https://github.com/RazvanBerbece) for the configuration)

## Running the full service composition
1. Run a freshly built full service composition (app, DBs, etc.) with the `docker compose up -d --remove-orphans --build` command.
2. Run the existing DB migrations with the `goose mysql ... up` command specified in the `bot-service/data` [README.md](./internal/bot-service/data/README.md).
    - This is required so the local development database is configured with all the necessary default data.   
3. Bring down all the services by running `docker compose down -v`.

# CI/CD
This project will employ CI/CD through the use of GitHub Actions and (probably?) Microsoft Azure. 

## CI
Continuous integration will be implemented through a workflow script which sets up a Go environment and then runs the internal logic tests on all pull request and pushes to main. The workflow file for CI can be seen in [test.yml](.github/workflows/test.yml).

## CD
Continuous deployment will be implemented through a workflow script which builds all the project artifacts and pushes them to Google Cloud on pushes to the main branch. The workflow file for CD can be seen in [deploy.yml](.github/workflows/deploy.yml).

Notes:
1. The production environment secret is base64 encoded using `openssl base64 -A -in .prod.env -out .env.out`
2. The production DB connection string is base64 encoded using `echo -n "CONN_STRING" | openssl base64`

# Contribution Guide
## Folder Structure
1. `cmd` folder -- contains multiple folders, each one representing a service making the `Aztebot` system. For example, the `bot-service` folder is the main entry point of the bot application (and it contains the Dockerfile associated with it) that starts the connection to Discord and actions on the various events emitted by it. More services can be added here by adding new folders (e.g.: `leveling-service/`) with `main.go` entry points and `Dockerfile`s which leverage these services.
2. `internal` folder -- has multiple folders, each one containing the internal logic that each service needs. Each service internals subfolder contains everything from data models and migration histories, to loggers, to handlers and contexts.
3. `pkg` folder -- contains util packages which could be leveraged across the entire project.
