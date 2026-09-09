package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/server"
)

type HealthHandler struct {
	Handler
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		Handler{server: s},
	}
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	start := time.Now()
	logger := middleware.GetLogger(c).With().
		Str("operation", "health_check").
		Logger()

	response := map[string]interface{}{
		"status":      "healthy",
		"timestamp":   time.Now().Unix(),
		"environment": h.server.Config.Primary.Env,
		"endpoint":    "/health",
		"checks":      make(map[string]interface{}),
	}

	checks := response["checks"].(map[string]interface{})
	isHealthy := true

	// Check database connectivity
	dbCtx, dbCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer dbCancel()

	dbStart := time.Now()
	if err := h.server.DB.Pool.Ping(dbCtx); err != nil {
		checks["database"] = map[string]interface{}{
			"status":        "unhealthy",
			"response_time": time.Since(dbStart).String(),
			"error":         err.Error(),
		}
		isHealthy = false
		logger.Error().Err(err).Dur("duration", time.Since(dbStart)).Msg("database health check failed")

		if h.server.LoggerService != nil && h.server.LoggerService.GetApplication() != nil {
			h.server.LoggerService.GetApplication().RecordCustomEvent(
				"HealthCheckError", map[string]interface{}{
					"check_type":       "database",
					"operation":        "health_check",
					"error_type":       "database_unhealthy",
					"response_time_ms": time.Since(dbStart).Milliseconds(),
					"error_message":    err.Error(),
				},
			)
		}
	} else {
		checks["database"] = map[string]interface{}{
			"status":        "healthy",
			"response_time": time.Since(dbStart).String(),
		}
		logger.Info().Msg("database health check passed")
	}

	// Check Redis connectivity — independent of database check
	if h.server.Redis != nil {
		redisCtx, redisCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer redisCancel()

		redisStart := time.Now()
		if err := h.server.Redis.Ping(redisCtx).Err(); err != nil {
			checks["redis"] = map[string]interface{}{
				"status":        "unhealthy",
				"response_time": time.Since(redisStart).String(),
				"error":         err.Error(),
			}
			isHealthy = false
			logger.Error().Err(err).Dur("response_time", time.Since(redisStart)).Msg("redis health check failed")

			if h.server.LoggerService != nil && h.server.LoggerService.GetApplication() != nil {
				h.server.LoggerService.GetApplication().RecordCustomEvent(
					"HealthCheckError", map[string]interface{}{
						"check_type":       "redis",
						"operation":        "health_check",
						"error_type":       "redis_unhealthy",
						"response_time_ms": time.Since(redisStart).Milliseconds(),
						"error_message":    err.Error(),
					})
			}
		} else {
			checks["redis"] = map[string]interface{}{
				"status":        "healthy",
				"response_time": time.Since(redisStart).String(),
			}
			logger.Info().Dur("response_time", time.Since(redisStart)).Msg("redis health check passed")
		}
	}

	// Set overall status
	if !isHealthy {
		response["status"] = "unhealthy"
		logger.Warn().Dur("duration", time.Since(start)).Msg("health check failed")

		if h.server.LoggerService != nil && h.server.LoggerService.GetApplication() != nil {
			h.server.LoggerService.GetApplication().RecordCustomEvent(
				"HealthCheckError", map[string]interface{}{
					"check_type":       "overall",
					"operation":        "health_check",
					"error_type":       "overall_unhealthy",
					"response_time_ms": time.Since(start).String(),
				})
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	logger.Info().Dur("duration", time.Since(start)).Msg("health check passed")
	c.JSON(http.StatusOK, response)
}
