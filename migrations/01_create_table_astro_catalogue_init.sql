-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DROP TABLE IF EXISTS astro_catalogue;

CREATE TABLE astro_catalogue (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL
);

ALTER TABLE astro_catalogue ADD CONSTRAINT primary_key_unique_pair PRIMARY KEY (name);

CREATE INDEX name_of_index ON astro_catalogue(lower(name));

-- +goose Down
