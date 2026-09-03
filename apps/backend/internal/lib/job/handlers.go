package job

import (
	"context"
	"encoding/json/v2"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/lib/email"
)

var emailClient *email.Client

func (job *JobService) InitHandlers(cfg *config.Config, logger *zerolog.Logger) {
	emailClient = email.NewClient(cfg, logger)
}

func (job *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	job.logger.Info().
		Str("type", "welcome_email").
		Str("to", payload.To).
		Msg("Processing welcome email task")

	err := emailClient.SendWelcomeEmail(payload.To, payload.FirstName)
	if err != nil {
		job.logger.Error().
			Str("type", "welcome_email").
			Str("to", payload.To).
			Err(err).
			Msg("Failed to send welcome email")

		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	job.logger.Info().
		Str("type", "welcome_email").
		Str("to", payload.To).
		Msg("Successfully sent welcome email")
	return nil
}
