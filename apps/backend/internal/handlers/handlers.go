package handlers

import (
	"github.com/xanity-07/openmat/internal/server"
	"github.com/xanity-07/openmat/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Users   *UserHandler
	Auth    *AuthHandler
}

func NewHandlers(s *server.Server, services *service.Service) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Users:   NewUserHandler(s, services.User),
		Auth:    NewAuthHandler(s, services.Auth),
	}
}
