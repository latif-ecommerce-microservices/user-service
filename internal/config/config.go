package config

import (
	"strings"

	"github.com/latif-ecommerce-microservices/user-service/pkg/envparser"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type Config struct {
	LogLevel    string `env:"LOG_LEVEL"`
	AppHTTPPort string `env:"APP_HTTP_PORT" default:"8080"`
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Host             string `env:"DB_HOST"`
	Port             string `env:"DB_PORT"`
	User             string `env:"DB_USERNAME"`
	Password         string `env:"DB_PASSWORD"`
	Name             string `env:"DB_NAME"`
	SSLMode          string `env:"DB_SSL_MODE"`
	MaxIdleConns     int    `env:"DB_MAX_IDLE_CONNS" default:"5"`
	MaxOpenConns     int    `env:"DB_MAX_OPEN_CONNS" default:"20"`
	EnableMigrations bool   `env:"ENABLE_MIGRATIONS" default:"true"`
}

func (c *Config) GetLogLevel() logging.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return logging.LevelDebug
	case "info":
		return logging.LevelInfo
	case "warn":
		return logging.LevelWarn
	case "error":
		return logging.LevelError
	case "fatal":
		return logging.LevelFatal
	default:
		return logging.LevelInfo
	}
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := envparser.LoadEnv(&cfg)
	if err != nil {
		return nil, nil
	}

	return &cfg, nil
}
