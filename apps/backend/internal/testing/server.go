package testing

import (
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/database"
	"github.com/xanity-07/openmat/internal/server"
)

// CreateTestServer creates a server instance for testing
func CreateTestServer(logger *zerolog.Logger, db *TestDB) *server.Server {
	// Set up observability config with defaults if not present
	if db.Config.Observability == nil {
		db.Config.Observability = &config.ObservabilityConfig{
			ServiceName: "tasker-test",
			Environment: "test",
			Logging: config.LoggingConfig{
				Level:  "info",
				Format: "json",
			},
			NewRelic: config.NewRelicConfig{
				LicenseKey:               "",    // Empty for tests
				AppLogForwardingEnabled:  false, // Disabled for tests
				DistributedTracerEnabled: false, // Disabled for tests
				DebugLogger:              false, // Disabled for tests
			},
			HealthChecks: config.HealthChecksConfig{
				Enabled: false,
			},
		}
	}

	testServer := &server.Server{
		Logger: logger,
		DB: &database.Database{
			Pool: db.Pool,
		},
		Config: db.Config,
	}

	return testServer
}
