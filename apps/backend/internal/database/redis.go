package database

import (
	"context"
	"fmt"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/loggerpkg"
)

func NewRedis(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerpkg.LoggerService) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.Address)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse redis address")
		return nil, fmt.Errorf("failed to parse redis address: %w", err)
	}

	rdb := redis.NewClient(opts)

	if loggerService != nil && loggerService.GetApplication() != nil {
		rdb.AddHook(nrredis.NewHook(rdb.Options()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error().Err(err).Msg("failed to connect to redis, continuing without Redis")
	} else {
		logger.Info().Str("component", "Redis").Msg("Connected to Redis successfully")
	}

	return rdb, nil
}
