-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"
(
    id       uuid NOT NULL PRIMARY KEY,
    login    text NOT NULL UNIQUE,
    password text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd