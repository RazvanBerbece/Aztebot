# azteca-discord
The ambitious and robust Discord bot which powers the <INSERT DISCORD NAME HERE ONCE WE KNOW IT> Discord community. Written in Go.

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
1. Run a built full service composition (app, DBs, etc.) with the `docker compose up -d --remove-orphans --build` command.
2. Bring down all the services by running `docker compose down`.

# CI/CD
This project will employ CI/CD through the use of GitHub Actions and (probably?) Microsoft Azure. 

## CI
Continuous integration will be implemented through a workflow script which sets up a Go environment and then runs the internal logic tests on all pull request and pushes to main.

## CD
The continuous deployment process has not been determined yet.

# Contribution Guide
## Folder Structure
1. `cmd` folder -- contains multiple folders, each one representing a service making the `azteca-discord` bot project. For example, the `bot-service` folder is the main entry point of the bot application (and it contains the Dockerfile associated with it) that starts the connection to Discord and actions on the various events emitted by it. More services can be added here by adding new folders (e.g.: `leveling-service/`) with `main.go` entry points and `Dockerfile`s which leverage these services.
2. `internal` folder -- has multiple folders, each one containing the internal logic that each service needs. Each service internals subfolder contains everything from data models and migration histories, to loggers, to handlers and contexts.
3. `pkg` folder -- contains util packages which could be leveraged across the entire project.