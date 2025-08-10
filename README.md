# Dota 2 Pro Games Database

A comprehensive, up-to-date database of professional Dota 2 matches, extracted from the [Stratz API](https://stratz.com/). This project enables researchers, analysts, and fans to explore pro-level Dota 2 data with ease.

## Features

- **Automated Data Ingestion:** Fetches and imports all major professional leagues (with prize pool filtering) and their matches.
- **Rich Schema:** Stores detailed data on leagues, teams, series, matches, players, heroes, items, and more.
- **Up-to-date Constants:** Keeps hero, item, and game version data current.
- **Flexible Storage:** Supports both local SQLite and remote [Turso (libsql)](https://turso.tech/) databases.
- **Weekly Database Releases:** Download the latest pre-built SQLite database from GitHub Releases.

## Download Latest Database

A fresh SQLite database is published weekly, containing all completed leagues from Feb 16, 2025 onward (with prize pool ≥ $1,000,000).

[Download latest sqlite database](https://github.com/dca123/dota-pro-db/releases/tag/latest)

## Setup

### 1. Requirements

- **Go** 1.22 or newer
- **Stratz API Key** ([get one here](https://stratz.com/))
- (Optional) **Turso credentials** for remote DB

### 2. Environment Variables

Set your API key (and Turso credentials if using remote DB):

```sh
export STRATZ_API_KEY="your_stratz_api_key"
# For Turso (libsql) remote DB:
export TURSO_DATABASE_URL="your_turso_db_url"
export TURSO_AUTH_TOKEN="your_turso_auth_token"
```

## Usage

### 1. Initialize Database with Constants

Fetches and inserts all heroes, items, and game versions into your database. This is a required first step before importing any match data.

```sh
cd cmd/sync_heroes_items
go run main.go
```

This command will:
- Connect to the Stratz API
- Fetch all hero data (ID, name, attributes)
- Fetch all item data (ID, name)
- Fetch all game version data
- Insert this information into your database

### 2. Import a Specific League

Imports all matches for a given league ID. Useful when you want to add data for a specific tournament.

```sh
cd cmd/sync_league
# Replace {LEAGUE_ID} with the desired league's ID
go run main.go {LEAGUE_ID}
```

This command will:
- Fetch all matches for the specified league
- Import complete match data including:
  - Teams and players
  - Series information
  - Match details (duration, winner, etc.)
  - Hero picks and bans
  - Player performance stats
  - Items purchased

### 3. Bulk Import All Pro Leagues

Fetches and imports all pro leagues (with prize pool ≥ $1,000,000) since Feb 16, 2025.

```sh
cd cmd/get_pro_leagues
# For local SQLite DB
go run main.go
# For remote Turso DB
go run main.go --turso
```

This command will:
- Fetch a list of all professional leagues meeting the criteria
- Check which leagues are already in your database
- Import all missing leagues and their matches
- Skip leagues that are already imported

### 4. Using the `--turso` Flag

To use a remote Turso database instead of local SQLite, add the `--turso` flag to any command:

```sh
go run main.go --turso
```

Make sure you've set the required environment variables (`TURSO_DATABASE_URL` and `TURSO_AUTH_TOKEN`).

## Database Schema

The database includes the following tables:

### Core Tables

- **heroes**: All Dota 2 heroes with their IDs, names, and primary attributes
- **items**: All Dota 2 items with their IDs and names
- **game_versions**: All Dota 2 game versions with release dates

### League and Match Data

- **leagues**: Professional tournaments and their details
- **teams**: Professional teams that participated in matches
- **team_players**: Professional players with their Steam account IDs
- **series**: Best-of-X series between teams
- **matches**: Individual games with detailed statistics
- **match_pick_bans**: Hero picks and bans for each match
- **match_players**: Detailed player performance for each match

### Schema Migrations

Database migrations are managed in the `migrations/` directory:
- `20250327125848-init.sql`: Initial schema with heroes and items tables
- `20250328233611-league.sql`: League, team, and match tables
- `20250612021121-add-hero-stats-image.sql`: Added hero attributes
- `20250617173732-add-game-versions-table.sql`: Game versions table

## Automated Updates

This project includes GitHub Actions workflows that automatically update the database:

### Weekly SQLite Database Update

A GitHub Action runs every Sunday to:
1. Fetch all new professional matches
2. Update the SQLite database
3. Publish the updated database as a GitHub Release

### Weekly Turso Database Update

A separate GitHub Action runs every Sunday to update the remote Turso database with the latest matches.

## License & Contributing

- Contributions are welcome! Feel free to open issues or PRs.
- This project is licensed under the MIT License.

## Acknowledgements

- [Stratz API](https://stratz.com/) for providing rich Dota 2 data.
- [Turso](https://turso.tech/) for remote database support.
- [SQLite](https://www.sqlite.org/) for the embedded database engine.
- [Go-SQLite3](https://github.com/mattn/go-sqlite3) for SQLite bindings for Go.
- [SQL-Migrate](https://github.com/rubenv/sql-migrate) for database migrations.
