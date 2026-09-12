package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/database"
	handlers2 "github.com/xanity-07/openmat/internal/handlers"
	"github.com/xanity-07/openmat/internal/loggerpkg"
	"github.com/xanity-07/openmat/internal/repository"
	"github.com/xanity-07/openmat/internal/router"
	"github.com/xanity-07/openmat/internal/server"
	"github.com/xanity-07/openmat/internal/service"
)

const DefaultContextTimeout = 10

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	gin.SetMode(gin.ReleaseMode)

	// Initialize New Relic logger service
	loggerService := loggerpkg.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	// Initialize our logger with service if provided
	log := loggerpkg.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env == "development" {
		if err = database.Migrate(context.Background(), cfg, &log); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	// Initialize Redis
	redisClient, err := database.NewRedis(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	// Initialize server
	srv, err := server.New(cfg, &log, loggerService, redisClient)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}

	// Initialize repositories, services, handlers
	repos := repository.NewRepositories(srv)
	services := service.NewService(srv, repos)
	handlers := handlers2.NewHandlers(srv, services)

	// Initialize router
	r := router.NewRouter(srv, handlers, cfg, repos.SessionRepo)

	// Set-up HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interruption signal to gracefully shut down the server
	// 30 second timeout allows for any currently ongoing request to get processed
	// but any new incoming requests will get dropped
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*DefaultContextTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown server")
	}
	stop()
	cancel()

	log.Info().Msg("Server exited properly")
}
