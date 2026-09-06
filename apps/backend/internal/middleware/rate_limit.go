package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/server"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	server *server.Server
}

func NewRateLimitMiddleware(s *server.Server) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		server: s,
	}
}

func (r *RateLimitMiddleware) RecordRateLimitHit(endpoint string) {
	if r.server.LoggerService != nil && r.server.LoggerService.GetApplication() != nil {
		r.server.LoggerService.GetApplication().RecordCustomEvent("RateLimitHit", map[string]interface{}{
			"endpoint": endpoint,
		})
	}
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func (r *RateLimitMiddleware) RateLimiter() gin.HandlerFunc {
	var (
		mu      sync.Mutex
		clients map[string]*client
	)

	// Periodically remove inactive clients
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()

			mu.Lock()
			for ip, cl := range clients {
				if now.Sub(cl.lastSeen) > time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()

		cl, exists := clients[ip]
		if !exists {
			cl = &client{
				limiter: rate.NewLimiter(10, 20),
			}
		}

		cl.lastSeen = time.Now()
		allowed := cl.limiter.Allow()
		mu.Unlock()

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, errs.Response[any]{
				Success: false,
				Status:  http.StatusTooManyRequests,
				Error: &errs.ErrorInfo{
					Code:    errs.MakeUpperCaseWithUnderscores(http.StatusText(http.StatusTooManyRequests)),
					Message: "rate limit exceeded",
				},
			})
		}
		c.Next()
	}
}
