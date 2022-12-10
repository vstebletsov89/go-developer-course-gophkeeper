package server

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
)

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled":
		return zerolog.Disabled
	}
	return zerolog.InfoLevel
}

func RunServer(cfg *config.Config) error {
	// init global logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// default log level is info
	zerolog.SetGlobalLevel(parseLogLevel(cfg.LogLevel))

	// debug config
	log.Debug().Msgf("%+v\n\n", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := connectDB(ctx, cfg.DatabaseDsn)
	if err != nil {
		log.Debug().Msgf("connectDB error: %s", err)
		return err
	}
	log.Info().Msg("Connection OK")
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Error().Msgf("Close connection error: %s", err)
		}
	}(conn, ctx)

	return nil
}

func connectDB(ctx context.Context, databaseURL string) (*pgx.Conn, error) {
	log.Debug().Msg("Connect to DB...")
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Error().Msgf("Failed to connect to database. Error: %v", err.Error())
		return nil, err
	}
	return conn, nil
}
