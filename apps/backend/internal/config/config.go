// Package config handles all the configuration logic for the backend and observability.
package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	// _ "github.com/joho/godotenv/autoload" // Auto loads all env variables from .env
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability" validate:"required"`
	Integration   IntegrationConfig    `koanf:"integration" validate:"required"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            string `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConn     int    `koanf:"max_open_conn" validate:"required"`
	MaxIdleConn     int    `koanf:"max_idle_conn" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
	TTLHours  int    `koanf:"ttl_hours" validate:"required"`
}

type IntegrationConfig struct {
	ResendAPIKey string `koanf:"resend_api_key"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")
	if os.Getenv("OPENMAT_ENV") != "docker" {
		if err := godotenv.Load(); err != nil {
			logger.Warn().Err(err).Msg("failed to load .env")
		}
	}

	err := k.Load(env.Provider(".", env.Opt{
		Prefix: "OPENMAT_",
		TransformFunc: func(k string, v string) (string, any) {
			k = strings.ToLower(strings.TrimPrefix(k, "OPENMAT_"))
			k = strings.ReplaceAll(k, "__", ".")
			return k, v
		},
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load initial env variables")
	}

	// Unmarshal into an instance of Config
	mainConfig := &Config{}
	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to unmarshal into main config")
	}

	// Validate the mainConfig struct
	validate := validator.New()
	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to validate main config struct")
	}

	// Check if observability was provided and if not give a default
	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	// Override Service Name and Environment
	mainConfig.Observability.ServiceName = "openmat"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	// Validate observability
	err = mainConfig.Observability.Validate()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to validate observability config")
	}

	return mainConfig, nil
}
