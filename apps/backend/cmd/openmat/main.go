package main

import (
	"fmt"

	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/database"
	"github.com/xanity-07/openmat/internal/lib/utils"
	"github.com/xanity-07/openmat/internal/loggerpkg"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load app configurations")
	}

	loggerService := loggerpkg.NewLoggerService(cfg.Observability)
	log := loggerpkg.NewLoggerWithService(cfg.Observability, loggerService)

	_, err = database.New(cfg, &log, loggerService)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize PostgreSQL")
	}

	_, err = database.NewRedis(cfg, &log, loggerService)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize Redis")
	}

	id := utils.GenerateID(11)

	fmt.Println(id)
}
