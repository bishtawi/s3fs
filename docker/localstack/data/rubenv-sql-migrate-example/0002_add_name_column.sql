-- +migrate Up
ALTER TABLE people ADD COLUMN name TEXT;

-- +migrate Down
ALTER TABLE people DROP COLUMN name;
