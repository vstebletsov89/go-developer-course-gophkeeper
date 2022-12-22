package testhelpers

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"os"
)

// IsGithubActions checks that tests are running from github actions.
func IsGithubActions() bool {
	ci := os.Getenv("CI") == "true"
	log.Debug().Msgf("GithubActions: %v", ci)
	return ci
}

// ConnectDBForMigration connects to postgres database.
func ConnectDBForMigration(databaseURL string) (*sql.DB, error) {
	log.Debug().Msg("connectDBForMigration...")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
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
