// Package server contains all the logic containing our main backend services
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/database"
	"github.com/xanity-07/openmat/internal/lib/job"
	"github.com/xanity-07/openmat/internal/loggerpkg"
)

type Server struct {
	Config        *config.Config
	DB            *database.Database
	Redis         *redis.Client
	Logger        *zerolog.Logger
	LoggerService *loggerpkg.LoggerService
	Job           *job.JobService
	httpServer    *http.Server
}

// New returns an initialized server struct with all the dependencies
func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerpkg.LoggerService, redis *redis.Client) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize background job services
	jobService := job.NewJobService(cfg, logger)
	jobService.InitHandlers(cfg, logger)

	// Start the job server
	if err := jobService.Start(); err != nil {
		return nil, err
	}

	server := &Server{
		Config:        cfg,
		DB:            db,
		Redis:         redis,
		Logger:        logger,
		LoggerService: loggerService,
		Job:           jobService,
	}

	// Start metrics collection
	// Runtime metrics are automatically collected by New Relic Go agent

	return server, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("environment", s.Config.Primary.Env).
		Msg("Starting server")

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.Redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	s.Logger.Info().Msg("Closing Redis connection")

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close PostgreSQL connection")
	}

	if s.Job != nil {
		s.Job.Stop()
	}

	return nil
}
