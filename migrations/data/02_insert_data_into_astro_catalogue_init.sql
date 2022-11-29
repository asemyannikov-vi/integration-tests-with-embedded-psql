-- +goose Up

INSERT INTO astro_catalogue (name) VALUES ('Mercury');
INSERT INTO astro_catalogue (name) VALUES ('Venus');
INSERT INTO astro_catalogue (name) VALUES ('Earth');

-- +goose Down
