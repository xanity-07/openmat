// Package job contains all our background job services
package job

import (
	"strings"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/config"
)

type JobService struct {
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

func NewJobService(cfg *config.Config, logger *zerolog.Logger) *JobService {
	redisAddr := strings.TrimPrefix(cfg.Redis.Address, "redis://")

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6, // Higher priority queue for important tasks
				"default":  3, // Default priority for most tasks
				"low":      1, // Low priority for non-urgent tasks
			},
		},
	)

	return &JobService{
		Client: client,
		server: server,
		logger: logger,
	}
}

func (job *JobService) Start() error {
	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, job.handleWelcomeEmailTask)

	job.logger.Info().Msg("Starting background job service")
	if err := job.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (job *JobService) Stop() {
	job.logger.Info().Msg("Stopping background job service")
	job.server.Shutdown()
	job.Client.Close()
}
