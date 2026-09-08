package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/handlers"
)

func registerSystemRoutes(r *gin.Engine, h *handlers.Handlers) {
	r.GET("status", h.Health.CheckHealth)

	r.Static("/static", "static")

	r.GET("/docs", h.OpenAPI.ServerOpenAPIUI)
}
