-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "data"
(
    id          uuid        NOT NULL PRIMARY KEY,
    user_id     uuid REFERENCES users (id) ON DELETE CASCADE,
    data_type   integer     NOT NULL,
    data_binary bytea       NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "data";
-- +goose StatementEnd