package postgres

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// RunMigrations runs goose migrations for current postgres db.
func RunMigrations(db *sql.DB) error {
	log.Info().Msg("Migration started...")

	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}

	log.Info().Msg("Migration done")
	return nil
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
