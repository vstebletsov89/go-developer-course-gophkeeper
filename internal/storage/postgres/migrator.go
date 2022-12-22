package postgres

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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
