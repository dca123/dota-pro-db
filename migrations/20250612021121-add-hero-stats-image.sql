
-- +migrate Up
ALTER TABLE heroes ADD COLUMN primary_attribute TEXT NOT NULL DEFAULT "";

-- +migrate Down
ALTER TABLE heroes DROP COLUMN primary_attribute;
