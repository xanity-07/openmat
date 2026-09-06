package middleware

import "github.com/xanity-07/openmat/internal/server"

type Middlewares struct {
	TracingMiddleware *TracingMiddleware
	ContextEnhancer   *ContextEnhancer
	Global            *GlobalMiddlewares
	RateLimit         *RateLimitMiddleware
}

func NewMiddlewares(s *server.Server) *Middlewares {
	nrApp := s.LoggerService.GetApplication()

	return &Middlewares{
		TracingMiddleware: NewTracingMiddleware(s, nrApp),
		ContextEnhancer:   NewContextEnhancer(s),
		Global:            NewGlobalMiddleware(s),
		RateLimit:         NewRateLimitMiddleware(s),
	}
}
