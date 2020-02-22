-- +migrate Up
CREATE TABLE people (id INT);

-- +migrate Down
DROP TABLE people;
