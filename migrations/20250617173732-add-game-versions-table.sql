
-- +migrate Up
CREATE TABLE game_versions(
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  asOfDateTime DATETIME NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS game_versions;
