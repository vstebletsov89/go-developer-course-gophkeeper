// Package config provides settings for server and client.
package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"sync"
)

// Config contains global settings of service.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	DatabaseDsn   string `env:"DATABASE_DSN" envDefault:"user=postgres dbname=postgres password=postgres host=localhost sslmode=disable" json:"database_dsn"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"debug"`
}

var once sync.Once

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

	log.Printf("%+v\n\n", cfg)
	return &cfg, nil
}
