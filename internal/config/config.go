// Package config provides settings for server and client.
package config

import (
	"flag"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config contains global settings of service.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"serverAddress"`
	DatabaseDsn   string `env:"DATABASE_DSN" envDefault:"user=postgres dbname=postgres password=postgres host=localhost sslmode=disable" json:"databaseDsn"` //nolint:lll
	LogLevel      string `env:"LOG_LEVEL" envDefault:"debug"`
}

var once sync.Once //nolint:gochecknoglobals

func (c *Config) readCommandLineArgs() {
	once.Do(func() {
		flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server and port to listen on")
		flag.StringVar(&c.DatabaseDsn, "d", c.DatabaseDsn, "database dsn")
		flag.StringVar(&c.LogLevel, "l", c.LogLevel, "log level")
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
