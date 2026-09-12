package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/handlers"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/repository"
)

func registerUserRoutes(r gin.IRouter, h *handlers.UserHandler, jwtSecret string, sessionRepo *repository.SessionRepository) {
	users := r.Group("/users")
	users.POST("", h.CreateUser())

	// Admin routes
	admin := users.Group("")
	admin.Use(middleware.RequireAuth(jwtSecret, sessionRepo))
	admin.Use(middleware.RequireRole(string(enums.ADMIN)))
	admin.GET("", h.GetUsers())
}
