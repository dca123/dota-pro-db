-- +migrate Up
CREATE TABLE heroes (
  id INTEGER PRIMARY KEY, 
  display_name TEXT NOT NULL,
  short_name TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +migrate StatementBegin
CREATE TABLE items (
  id INTEGER PRIMARY KEY,
  display_name TEXT NOT NULL,
  short_name TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER update_hero_timestamp
AFTER UPDATE ON heroes
BEGIN
  UPDATE heroes
  SET updated_at = CURRENT_TIMESTAMP
  WHERE id = NEW.id;
END;

CREATE TRIGGER update_item_timestamp
AFTER UPDATE ON items
BEGIN
  UPDATE items
  SET updated_at = CURRENT_TIMESTAMP
  WHERE id = NEW.id;
END;
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS update_hero_timestamp;
DROP TRIGGER IF EXISTS update_item_timestamp;
DROP TABLE IF EXISTS heroes;
DROP TABLE IF EXISTS items;
