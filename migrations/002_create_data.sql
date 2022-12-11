-- This is a migration of data table.

CREATE TABLE IF NOT EXISTS data (
    id          uuid        NOT NULL PRIMARY KEY,
    user_id     uuid        REFERENCES users(id) ON DELETE CASCADE,
    data_type   text        NOT NULL,
    data_binary bytea       NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----

drop table data;
