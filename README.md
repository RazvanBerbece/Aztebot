# AzteBot
The powerful Discord bot which powers the OTA (Ordinul Templierilor Azteci) Discord community. Written in Go.

## Composing services
### Core
- `bot-service` (Handles Discord interactions like new messages, slash commands, join events, reaction adding or removing, etc.)
### Dependencies
- `mysql-db` (Containerised MySQL instance for local development)

### External
- `aztemusic-service` (Standalone music orchestrator bot application - see development here [AzteMusic](https://github.com/AzteBot-Developments/AzteMusic))

# How to Run
## Prerequisites
In order to run the application, a few prerequisites must be met.
1. Have the repository cloned locally.
2. Have Docker installed.
3. Have Make installed.
4. Have a fully-configured `.env` file saved in the root of the repository. (contact [@RazvanBerbece](https://github.com/RazvanBerbece) for the configuration)
5. Additionally, for full local development capabilities, have the [Aztebot-Infrastructure](https://github.com/RazvanBerbece/Aztebot-Infrastructure) repository cloned locally in a folder which also contains the `Aztebot` repository - i.e. Folder `Project` should contain both the `Aztebot` and the `Aztebot-Infrastructure` repository folders.

## Running the full service composition
1. Run a freshly built full service composition (app, DBs, etc.) with the `make run-all` command.
    - This is required so the local development database is configured with all the necessary default data.   
2. Run migrations locally by executing the following commands from the root of this repository
    - To execute a dryrun and double-check the to-be-applied migrations: `make migrate-up-dry` 
    - To apply the migrations `make migrate-up`
3. Bring down all the services by running `make down`.

# CI/CD
This project will employ CI/CD through the use of GitHub Actions and Google Cloud. 

## CI
Continuous integration will be implemented through a workflow script which sets up a Go environment and then runs the internal logic tests on all pull request and pushes to main. The workflow file for the AzteBot CI can be seen in [test.yml](.github/workflows/test.yml).

## CD
Continuous deployment is implemented through a workflow script which builds all the project artifacts and uploads them to Google Cloud Artifact Registry on pushes to the main branch. Additionally, a GKE pod is created with the new container image and ultimately executed upstream to run the apps. The workflow file for the AzteBot CD can be seen in [deploy.yml](.github/workflows/deploy.yml).

Notes:
- The production environment file is base64 encoded using `make update-env` and decoded accordingly in the Actions workflows.

# Contribution Guide
## Folder Structure
1. `cmd` folder -- contains multiple folders, each one representing a service making the `Aztebot` system. For example, the `bot-service` folder is the main entry point of the bot application (and it contains the Dockerfile associated with it) that starts the connection to Discord and actions on the various events emitted by it. More services can be added here by adding new folders (e.g.: `leveling-service/`) with `main.go` entry points and `Dockerfile`s which leverage these services.
2. `internal` folder -- has multiple folders, each one containing the internal logic that each service needs. Each service internals subfolder contains everything from data models and migration histories, to loggers, to handlers and contexts.
3. `pkg` folder -- contains util packages which are leveraged across the entire project.
