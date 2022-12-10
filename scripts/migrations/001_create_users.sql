-- This is a migration of users table.

CREATE TABLE IF NOT EXISTS users (
    id       uuid NOT NULL PRIMARY KEY,
    login    text NOT NULL UNIQUE,
    password text NOT NULL
);

---- create above / drop below ----

drop table users;
