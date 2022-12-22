// Package config provides settings for server and client.
package config

import (
	"flag"
	"github.com/rs/zerolog"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config contains global settings of service.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"serverAddress"`
	DatabaseDsn     string `env:"DATABASE_DSN" envDefault:"user=postgres dbname=postgres password=postgres host=localhost sslmode=disable" json:"databaseDsn"` //nolint:lll
	JwtSecretKey    string `env:"JWT_SECRET" envDefault:"secret_key" json:"jwtSecretKey"`
	EnableTLS       bool   `env:"ENABLE_TLS" envDefault:"false" json:"enableTLS"`
	EnableMigration bool   `env:"ENABLE_MIGRATION" envDefault:"false" json:"enableMigration"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"debug"`
}

var once sync.Once //nolint:gochecknoglobals

func (c *Config) readCommandLineArgs() {
	once.Do(func() {
		flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server and port to listen on")
		flag.StringVar(&c.DatabaseDsn, "d", c.DatabaseDsn, "database dsn")
		flag.StringVar(&c.LogLevel, "l", c.LogLevel, "log level")
		flag.StringVar(&c.JwtSecretKey, "j", c.JwtSecretKey, "jwt secret key")
		flag.BoolVar(&c.EnableTLS, "s", c.EnableTLS, "enable secure mode")
		flag.BoolVar(&c.EnableMigration, "m", c.EnableMigration, "enable database migration")
		flag.Parse()
	})
}

// ReadConfig merges settings from environment and command line arguments.
func ReadConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}
	cfg.readCommandLineArgs()

	return &cfg, nil
}

// ParseLogLevel is a helper function to parse log level for zerolog.
func ParseLogLevel(level string) zerolog.Level {
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
