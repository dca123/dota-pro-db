-- +migrate Up
CREATE TABLE leagues (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE teams (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  tag TEXT NOT NULL
);

CREATE TABLE series (
  id INTEGER PRIMARY KEY,
  type TEXT NOT NULL,
  team_one_win_count INTEGER NOT NULL,
  team_two_win_count INTEGER NOT NULL,

  winning_team_id INTEGER,
  team_one_id INTEGER NOT NULL,
  team_two_id INTEGER NOT NULL,

  FOREIGN KEY (winning_team_id) REFERENCES teams (id),
  FOREIGN KEY (team_one_id) REFERENCES teams (id),
  FOREIGN KEY (team_two_id) REFERENCES teams (id)
);

CREATE TABLE matches (
  id INTEGER PRIMARY KEY,
  did_radiant_win INTEGER NOT NULL,
  duration_seconds INTEGER NOT NULL,
  start_date_time INTEGER NOT NULL,
  end_date_time INTEGER NOT NULL,
  tower_status_radiant INTEGER NOT NULL,
  tower_status_dire INTEGER NOT NULL,
  barracks_status_radiant INTEGER NOT NULL,
  barracks_status_dire INTEGER NOT NULL,
  first_blood_time INTEGER NOT NULL,
  lobby_type TEXT NOT NULL,
  game_mode TEXT NOT NULL,
  game_version_id INTEGER NOT NULL,
  radiant_networth_leads TEXT NOT NULL,
  radiant_experience_leads TEXT NOT NULL,
  analysis_outcome TEXT NOT NULL,

  league_id INTEGER NOT NULL,
  series_id INTEGER NOT NULL,
  radiant_team_id INTEGER NOT NULL,
  dire_team_id INTEGER NOT NULL,

  FOREIGN KEY (league_id) REFERENCES leagues (id),
  FOREIGN KEY (series_id) REFERENCES series (id),
  FOREIGN KEY (radiant_team_id) REFERENCES teams (id),
  FOREIGN KEY (dire_team_id) REFERENCES teams (id)
);
CREATE TABLE match_pick_bans (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  is_pick INTEGER NOT NULL,
  pick_order INTEGER NOT NULL,
  is_radiant INTEGER NOT NULL,

  match_id INTEGER NOT NULL,
  hero_id INTEGER NOT NULL,

  FOREIGN KEY (match_id) REFERENCES matches (id),
  FOREIGN KEY (hero_id) REFERENCES heroes (id)
);
CREATE TABLE match_players (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  is_radiant INTEGER NOT NULL,
  is_victory INTEGER NOT NULL,
  kills INTEGER NOT NULL,
  deaths INTEGER NOT NULL,
  assists INTEGER NOT NULL,
  num_last_hits INTEGER NOT NULL,
  num_denies INTEGER NOT NULL,
  gold_per_min INTEGER NOT NULL,
  networth INTEGER NOT NULL,
  exp_per_min INTEGER NOT NULL,
  level INTEGER NOT NULL,
  gold_spent INTEGER NOT NULL,
  hero_damage INTEGER NOT NULL,
  tower_damage INTEGER NOT NULL,
  hero_healing INTEGER NOT NULL,
  is_random INTEGER NOT NULL,
  lane TEXT NOT NULL,
  position TEXT NOT NULL,
  role TEXT NOT NULL,
  invisible_seconds INTEGER NOT NULL,
  dota_plus_hero_level INTEGER,

  match_id INTEGER NOT NULL,
  steam_account_id TEXT NOT NULL,
  hero_id INTEGER NOT NULL,
  item_0_id INTEGER,
  item_1_id INTEGER,
  item_2_id INTEGER,
  item_3_id INTEGER,
  item_4_id INTEGER,
  item_5_id INTEGER,
  backpack_0_id INTEGER,
  backpack_1_id INTEGER,
  backpack_2_id INTEGER,
  neutral_0_id INTEGER,

  FOREIGN KEY (match_id) REFERENCES matches (id),
  FOREIGN KEY (steam_account_id) REFERENCES team_players (steam_account_id),
  FOREIGN KEY (hero_id) REFERENCES heroes (id),
  FOREIGN KEY (item_0_id) REFERENCES items (id),
  FOREIGN KEY (item_1_id) REFERENCES items (id),
  FOREIGN KEY (item_2_id) REFERENCES items (id),
  FOREIGN KEY (item_3_id) REFERENCES items (id),
  FOREIGN KEY (item_4_id) REFERENCES items (id),
  FOREIGN KEY (item_5_id) REFERENCES items (id),
  FOREIGN KEY (backpack_0_id) REFERENCES items (id),
  FOREIGN KEY (backpack_1_id) REFERENCES items (id),
  FOREIGN KEY (backpack_2_id) REFERENCES items (id),
  FOREIGN KEY (neutral_0_id) REFERENCES items (id)
);

CREATE TABLE team_players (
  steam_account_id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS leagues;
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS match_pick_bans;
DROP TABLE IF EXISTS match_players;
DROP TABLE IF EXISTS series;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS team_players;
PRAGMA foreign_keys = OFF;
