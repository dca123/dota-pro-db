
-- +migrate Up
ALTER TABLE leagues ADD COLUMN liquipedia_page_name TEXT DEFAULT '';
CREATE INDEX idx_leagues_liquipedia_page_name ON leagues(liquipedia_page_name);

-- +migrate Down
DROP INDEX IF EXISTS idx_leagues_liquipedia_page_name;
ALTER TABLE leagues DROP COLUMN liquipedia_page_name;
