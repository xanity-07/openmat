// Package database contains all the logic related to connecting to a PostgreSQL connection pool Redis client and handles migrations
package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/newrelic/go-agent/v3/integrations/nrpgx5"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/loggerpkg"
)

const (
	DatabasePingTimeout = 10
)

type Database struct {
	Pool *pgxpool.Pool
	log  *zerolog.Logger
}

type multiTracer struct {
	tracers []any
}

func (mt *multiTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryStartData(context.Context, *pgx.Conn, pgx.TraceQueryStartData) context.Context
		}); ok {
			ctx = t.TraceQueryStartData(ctx, conn, data)
		}
	}
	return ctx
}

func (mt *multiTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryEndData(context.Context, *pgx.Conn, pgx.TraceQueryEndData)
		}); ok {
			t.TraceQueryEndData(ctx, conn, data)
		}
	}
}

// New creates a PostgreSQL connection pool
func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerpkg.LoggerService) (*Database, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, cfg.Database.Port)

	// URL-encoded password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf(
		"postgres://%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	// If we have a New Relic Application already, forward the tracer part of our pgxPoolConfig
	// from our database driver and add New Relic PostgreSQL instrumentation
	if loggerService != nil && loggerService.GetApplication() != nil {
		pgxPoolConfig.ConnConfig.Tracer = nrpgx5.NewTracer()
	}

	if cfg.Primary.Env == "local" {
		globalLevel := logger.GetLevel()
		pgxLogger := loggerpkg.NewPgxLogger(globalLevel)
		// Chain tracer - New Relic first, then local logging
		if pgxPoolConfig.ConnConfig.Tracer != nil {
			// If New Relic tracer exists, create multi-tracer
			localTracer := tracelog.TraceLog{
				Logger:   pgxzero.NewLogger(pgxLogger),
				LogLevel: tracelog.LogLevel(loggerpkg.GetPgxTraceLogLevel(globalLevel)),
			}

			pgxPoolConfig.ConnConfig.Tracer = &multiTracer{
				tracers: []any{pgxPoolConfig.ConnConfig.Tracer, &localTracer},
			}
		} else {
			pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
				Logger: pgxzero.NewLogger(pgxLogger),
			}
		}
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	database := &Database{
		Pool: pool,
		log:  logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL")
	}

	logger.Info().Str("component", "database").Msg("Connected to PostgreSQL successfully")
	return database, nil
}

// Close closes all connections in the pool and rejects future Pool.Acquire calls. Blocks until all connections are returned to pool and closed.
func (db *Database) Close() error {
	db.log.Info().Msg("closing database connection pool")
	db.Pool.Close()
	return nil
}
