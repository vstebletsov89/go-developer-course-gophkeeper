package testhelpers

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const postgreSQLTables = `
CREATE TABLE IF NOT EXISTS users (
    id       uuid NOT NULL PRIMARY KEY,
    login    text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS data (
    id          uuid        NOT NULL PRIMARY KEY,
    user_id     uuid        REFERENCES users(id) ON DELETE CASCADE,
    data_type   integer     NOT NULL,
    data_binary bytea       NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

// MigrateTables migrates all required tables for gophkeeper service.
func MigrateTables(pool *pgxpool.Pool) error {
	log.Debug().Msg("Migration started..")
	_, err := pool.Exec(context.Background(), postgreSQLTables)
	if err != nil {
		return err
	}
	log.Debug().Msg("Migration done")
	return nil
}

// ConnectDB connects to postgres database.
func ConnectDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	log.Debug().Msg("Connect to DB...")
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Error().Msgf("Failed to connect to database. Error: %v", err.Error())
		return nil, err
	}
	return pool, nil
}
