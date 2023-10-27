# azteca-discord
The ambitious and robust Discord bot which powers the <INSERT DISCORD NAME HERE> Discord community. Written in Go.

# How to Run
## Prerequisites
In order to run the application, a few prerequisites must be met.
1. Have the repository cloned locally.
2. Have the bot join a server and give it all the right permissions.
2. Have Docker installed.
3. Have a fully-configured `.env` file saved in the root of the repository. (contact @RazvanBerbece for the configuration)

## Running the full service composition
1. Run a built full service composition (app, DBs, etc.) with the `docker compose up -d --remove-orphans --build` command.
2. Bring down all the services by running `docker compose down`.

# CI/CD

## CI

## CD

# Contribution Guide
## Folder Structure
1. `cmd` folder -- contains multiple folders, each one representing a service making the `azteca-discord` bot project. For example, the `bot-service` folder is the main entry point of the bot application (and it contains the Dockerfile associated with it) that starts the connection to Discord and actions on the various events emitted by it. More services can be added here by adding new folders (e.g.: `leveling-service/`) with `main.go` entry points and `Dockerfile`s which leverage these services.
2. `internal` folder -- has multiple folders, each one containing the internal logic that each service needs. Each service internals subfolder contains everything from data models and migration histories, to loggers, to handlers and contexts.
3. `pkg` folder -- contains util packages which could be leveraged across the entire project.