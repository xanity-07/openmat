package config

import (
	"fmt"
	"time"
)

type ObservabilityConfig struct {
	ServiceName  string             `koanf:"service_name" validate:"required"`
	Environment  string             `koanf:"environment" validate:"required"`
	Logging      LoggingConfig      `koanf:"logging" validate:"required"`
	NewRelic     NewRelicConfig     `koanf:"new_relic" validate:"required"`
	HealthChecks HealthChecksConfig `koanf:"health_checks" validate:"required"`
}

type LoggingConfig struct {
	Level  string `koanf:"level" validate:"required"`
	Format string `koanf:"format" validate:"required"`
}

type NewRelicConfig struct {
	LicenseKey               string `koanf:"license_key"`
	AppLogForwardingEnabled  bool   `koanf:"app_log_forwarding_enabled"`
	DistributedTracerEnabled bool   `koanf:"distributed_tracer_enabled"`
	DebugLogger              bool   `koanf:"debug_logger"`
}

type HealthChecksConfig struct {
	Enabled  bool          `koanf:"enabled"`
	Timeout  time.Duration `koanf:"timeout" validate:"min=1s"`
	Interval time.Duration `koanf:"interval" validate:"min=1s"`
	Checks   []string      `koanf:"checks"`
}

func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		ServiceName: "openmat",
		Environment: "development",
		Logging: LoggingConfig{
			Level:  "debug",
			Format: "json",
		},
		NewRelic: NewRelicConfig{
			LicenseKey:               "",
			AppLogForwardingEnabled:  true,
			DistributedTracerEnabled: true,
			DebugLogger:              false,
		},
		HealthChecks: HealthChecksConfig{
			Enabled:  false,
			Timeout:  30,
			Interval: 5,
			Checks:   []string{"database", "redis"},
		},
	}
}

func (c *ObservabilityConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}

	validLevel := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}

	if !validLevel[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", c.Logging.Level)
	}

	return nil
}

func (c *ObservabilityConfig) GetLogLevel() string {
	switch c.Logging.Level {
	case "production":
		if c.Logging.Level == "" {
			return "info"
		}
	case "development":
		if c.Logging.Level == "" {
			return "debug"
		}
	}
	return c.Logging.Level
}

func (c *ObservabilityConfig) IsProduction() bool {
	return c.Logging.Level == "production"
}
