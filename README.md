# Dota 2 Pro games database

## Usage

- Set your API key
```
export STRATZ_API_KEY=""
```
### DB setup
- Go to the director `cmd/sync_heroes_items`
- Run the executable via `./sync_heroes_items`. 
- This will initialize the db for the first time with the required constants (heroes and items)

### Import Leagues
- Go to the director `cmd/sync_league`
- Run the executable via `./sync_league {LEAGUE_ID}`
- The db will be created in the `test.db` file in the root directory


