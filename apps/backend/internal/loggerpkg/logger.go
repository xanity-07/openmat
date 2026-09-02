// Package loggerpkg contains all the logic related to logging and New Relic observability.
package loggerpkg

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/xanity-07/openmat/internal/config"
)

type LoggerService struct {
	nrApp *newrelic.Application
}

func NewLoggerService(cfg *config.ObservabilityConfig) *LoggerService {
	service := &LoggerService{}

	// Checking if New Relic service is provided if not continue without it
	if cfg.NewRelic.LicenseKey == "" {
		return service
	}

	// Configuring New Relic based on env variables
	var configOptions []newrelic.ConfigOption
	configOptions = append(
		configOptions,
		newrelic.ConfigAppName(cfg.ServiceName),
		newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.AppLogForwardingEnabled),
		newrelic.ConfigDistributedTracerEnabled(cfg.NewRelic.DistributedTracerEnabled),
	)

	// Explicitly add DebugLogger only if experiencing issues with New Relic
	if cfg.NewRelic.DebugLogger {
		configOptions = append(configOptions, newrelic.ConfigDebugLogger(os.Stdout))
	}

	// Initializing a instance of a New Relic application if fails return empty service
	app, err := newrelic.NewApplication(configOptions...)
	if err != nil {
		return service
	}

	// Assigning the New Relic application to the service
	service.nrApp = app

	fmt.Printf("New Relic application initialized for service: %s\n", cfg.ServiceName)
	return service
}

// Shutdown terminates our New Relic application
func (ls *LoggerService) Shutdown() {
	if ls.nrApp != nil {
		ls.nrApp.Shutdown(time.Second * 10)
	}
}

// GetApplication returns the instance of our New Relic application
func (ls *LoggerService) GetApplication() *newrelic.Application {
	return ls.nrApp
}

// NewLogger creates a new Logger with specified log level
func NewLogger(level string, isProd bool) zerolog.Logger {
	return NewLoggerWithService(&config.ObservabilityConfig{
		Logging: config.LoggingConfig{
			Level: level,
		},
		Environment: func() string {
			if isProd {
				return "production"
			}
			return "development"
		}(),
	}, nil)
}

// NewLoggerWithService creates a Logger with full config and our New Relic LoggerService
func NewLoggerWithService(cfg *config.ObservabilityConfig, loggerService *LoggerService) zerolog.Logger {
	level := cfg.GetLogLevel()
	var logLevel zerolog.Level

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Don't set global level - let each logger have its own level
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Initializing our writer
	var writer io.Writer

	// This block of code will deal with whether we should forward our logs to New Relic or not
	var baseWriter io.Writer
	if cfg.IsProduction() && cfg.Logging.Format == "json" {
		// Production structured logs to stdout
		baseWriter = os.Stdout
		// Wrap with New Relic zerologWriter for log forwarding in production
		if loggerService != nil && loggerService.GetApplication() != nil {
			nrWriter := zerologWriter.New(baseWriter, loggerService.nrApp)
			writer = nrWriter
		} else {
			writer = baseWriter
		}
	} else {
		// Development mode - use console writer
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05 ",
		}
	}

	// Initialize the logger with logLevel and writer
	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("environment", cfg.Environment).
		Logger()

	if cfg.IsProduction() {
		logger.With().Stack().Logger()
	}

	return logger
}

// NewLoggerWithConfig creates a Logger with full config
func NewLoggerWithConfig(cfg *config.ObservabilityConfig) zerolog.Logger {
	return NewLoggerWithService(cfg, nil)
}

// WithTraceContext adds New Relic transaction context to logger
func WithTraceContext(logger zerolog.Logger, txn *newrelic.Transaction) zerolog.Logger {
	if txn == nil {
		return logger
	}

	// Get trace metadata from New Relic transaction
	metadata := txn.GetTraceMetadata()

	return logger.With().
		Str("trace.id", metadata.TraceID).
		Str("span.id", metadata.SpanID).
		Logger()
}

// NewPgxLogger creates a database Logger
func NewPgxLogger(level zerolog.Level) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatFieldValue: func(i interface{}) string {
			switch v := i.(type) {
			case string:
				// Clean and format long SQL statements for better readability
				if len(v) > 200 {
					return v[:200] + "..."
				}
				return v
			case []byte:
				var obj interface{}
				if err := json.Unmarshal(v, obj); err != nil {
					pretty, _ := json.Marshal(obj)
					return string(pretty)
				}
				return string(v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}

	return zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("component", "database").
		Logger()
}

// GetPgxTraceLogLevel converts a zerolog level into a pgx tracelog level
func GetPgxTraceLogLevel(level zerolog.Level) int {
	switch level {
	case zerolog.DebugLevel:
		return 6 // tracelog.LogLevelDebug
	case zerolog.InfoLevel:
		return 4 // tracelog.LogLevelInfo
	case zerolog.WarnLevel:
		return 3 // tracelog.LogLevelWarn
	case zerolog.ErrorLevel:
		return 2 // tracelog.LogLevelError
	default:
		return 0 // tracelog.LogLevelNone
	}
}
