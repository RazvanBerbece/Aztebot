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
5. Additionally, for full local development capabilities and to run the database migrations on the development machine, have the [Aztebot-Infrastructure](https://github.com/RazvanBerbece/Aztebot-Infrastructure) repository cloned locally in a folder which also contains the `Aztebot` repository (**For example**, the folder `Project` should contain both the `Aztebot` and the `Aztebot-Infrastructure` repository folders) 

_Note:_ At the moment, the Infrastructure submodule has to be updated when there are changes in the remote repository (e.g. a new migration file).

## Running the full service composition
1. Run a freshly built full service composition (app, DBs, etc.) with the `make up` command.
    - This is required so the local development database is configured with all the necessary default data.   
2. Run migrations locally by executing the following commands from the root of this repository (_requires [Aztebot-Infrastructure](https://github.com/RazvanBerbece/Aztebot-Infrastructure) as described in prerequisite #5_)
    - To execute a dryrun and double-check the to-be-applied migrations: `make migrate-up-dry` 
    - To apply the migrations `make migrate-up`

To bring down all the services, one can do so by running `make down`.

# CI/CD
This project will employ CI/CD through the use of GitHub Actions and Google Cloud. 

## CI
Continuous integration is implemented through a workflow script which sets up a containerised service composition containing the Go environment and other dependencies (MySql, etc.) and then runs the internal logic tests on all pull request and pushes to main. The workflow file for the AzteBot CI can be seen in [test.yml](.github/workflows/test.yml).

## CD
Continuous deployment is implemented through a workflow script which builds all the project artifacts and uploads them to Google Cloud Artifact Registry on pushes to the main branch. Additionally, a GKE pod is created with the new container image and ultimately executed upstream to run the apps. The workflow file for the AzteBot CD can be seen in [deploy.yml](.github/workflows/deploy.yml).

Notes:
- The production environment file is base64 encoded using `make update-envs` and decoded accordingly in the Actions workflows.

# Contribution Guide
## Folder Structure
1. `cmd` folder -- contains the main entry point of the bot application (and it also contains the Dockerfile associated with it) that starts the connection to Discord and actions on the various events emitted by it.
2. `internal` folder -- has multiple folders, each one containing the various logic components that the Aztebot service needs including data models, interfaces, handlers and DB contexts.
3. `pkg` folder -- contains util packages which are leveraged across the entire project.
