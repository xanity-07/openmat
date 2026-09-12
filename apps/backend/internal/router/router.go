package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/config"
	"github.com/xanity-07/openmat/internal/handlers"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/repository"
	v1 "github.com/xanity-07/openmat/internal/router/v1"
	"github.com/xanity-07/openmat/internal/server"
)

func NewRouter(s *server.Server, h *handlers.Handlers, cfg *config.Config, sessionRepo *repository.SessionRepository) *gin.Engine {
	mw := middleware.NewMiddlewares(s)

	router := gin.New()

	// Global middleware registration
	router.Use(
		//mw.RateLimit.RateLimiter(),
		mw.Global.Recover(),
		mw.Global.CORS(),
		mw.Global.Secure(),
		middleware.RequestID(),
		mw.Tracing.NewRelicMiddleware(),
		mw.Tracing.EnhanceTracing(),
		mw.ContextEnhancer.EnhanceContext(),
		mw.Global.RequestLogger(),
	)

	// register system routes
	registerSystemRoutes(router, h)

	// register v1 routes
	v1Router := router.Group("/api/v1")
	v1.RegisterV1Routes(v1Router, h, cfg.Auth.SecretKey, sessionRepo)

	return router
}
