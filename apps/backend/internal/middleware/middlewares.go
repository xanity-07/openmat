package middleware

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/xanity-07/openmat/internal/server"
)

type Middlewares struct {
	Tracing         *TracingMiddleware
	ContextEnhancer *ContextEnhancer
	Global          *GlobalMiddlewares
	RateLimit       *RateLimitMiddleware
}

func NewMiddlewares(s *server.Server) *Middlewares {
	var nrApp *newrelic.Application

	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}

	return &Middlewares{
		Tracing:         NewTracingMiddleware(s, nrApp),
		ContextEnhancer: NewContextEnhancer(s),
		Global:          NewGlobalMiddleware(s),
		RateLimit:       NewRateLimitMiddleware(s),
	}
}
