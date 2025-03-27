
-- +migrate Up
CREATE TABLE heroes (id INT PRIMARY KEY, display_name TEXT NOT NULL, short_name TEXT NOT NULL);
CREATE TABLE items (id INT PRIMARY KEY, display_name TEXT NOT NULL, short_name TEXT NOT NULL);

-- +migrate Down
DROP TABLE IF EXISTS heroes;
DROP TABLE IF EXISTS items;
