# Dota 2 Pro games database

## Download sqlite file
This file is updated each week. It includes all completed leagues from Feb 16, 2025 till the date the release was created with a prize pool >= $1,000,000.
[Release](https://github.com/dca123/dota-pro-db/releases/tag/latest)

## Usage

- Set your API key
```
export STRATZ_API_KEY=""
```
### DB setup
- Go to the director `cmd/sync_heroes_items`
- Run the program via `./go run main.go`. 
- This will initialize the db for the first time with the required constants (heroes and items)

### Import Leagues
- Go to the director `cmd/sync_league`
- Run the program via `./go run main.go {League_Id}`. 
- The db will be created in the `test.db` file in the root directory


