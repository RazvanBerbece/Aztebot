# Data
This folder contains data models and DB contexts leveraged by the bot. It also contains the `migrations` folder which contains migration SQL scripts.

## How to Run a Migration

### Locally
1. `cd` to the `migrations` folder and run all subsequent terminal commands in there.
2. Create a new migration SQL script using the `goose create <migration_name> sql` command.
3. Edit the generated SQL script to contain the desired DB schema changes.
4. Then run the `goose mysql "{root_username}:{root_password}@tcp(host:port)/{database_name}-{environment}" up` command to apply the migrations. (where the `{...}` constructs need to be replaced with the right details from the `.env` file).

### In the Cloud
This will probably work in a similar way to running the migrations locally, however need to find a more straightforward way to do this in the Cloud. To consider.