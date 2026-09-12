// internal/router/v1/router.go
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/handlers"
	"github.com/xanity-07/openmat/internal/repository"
)

func RegisterV1Routes(router gin.IRouter, h *handlers.Handlers, jwtSecret string, sessionRepo *repository.SessionRepository) {
	// Register User Routes
	registerUserRoutes(router, h.Users, jwtSecret, sessionRepo)

	// Register Auth Routes
	registerAuthRoutes(router, h.Auth, jwtSecret, sessionRepo)
}
