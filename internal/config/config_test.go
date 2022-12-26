package config

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name: "read config with defaults",
			want: &Config{
				ServerAddress:   "localhost:8080",
				DatabaseDsn:     "user=postgres dbname=postgres password=postgres host=localhost sslmode=disable",
				JwtSecretKey:    "secret_key",
				LogLevel:        "debug",
				EnableTLS:       false,
				EnableMigration: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{
			name:  "trace level",
			level: "trace",
			want:  zerolog.TraceLevel,
		},
		{
			name:  "debug level",
			level: "debug",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "info level",
			level: "info",
			want:  zerolog.InfoLevel,
		},
		{
			name:  "warn level",
			level: "warn",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "error level",
			level: "error",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "fatal level",
			level: "fatal",
			want:  zerolog.FatalLevel,
		},
		{
			name:  "panic level",
			level: "panic",
			want:  zerolog.PanicLevel,
		},
		{
			name:  "disabled level",
			level: "disabled",
			want:  zerolog.Disabled,
		},
		{
			name:  "default level",
			level: "default",
			want:  zerolog.InfoLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ParseLogLevel(tt.level), "ParseLogLevel(%v)", tt.level)
		})
	}
}
