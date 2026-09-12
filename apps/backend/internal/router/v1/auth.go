// internal/router/v1/auth.go
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/handlers"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/repository"
)

func registerAuthRoutes(r gin.IRouter, h *handlers.AuthHandler, jwtSecret string, sessionRepo *repository.SessionRepository) {
	auth := r.Group("/auth")
	auth.POST("/login", h.Login())

	// Grouped off `auth`, not off the top-level router `r` — logout needs
	// to land at /auth/logout, not /logout.
	protected := auth.Group("")
	protected.Use(middleware.RequireAuth(jwtSecret, sessionRepo))
	protected.POST("/logout", h.Logout())
}
