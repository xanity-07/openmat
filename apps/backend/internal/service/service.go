package service

import (
	"github.com/xanity-07/openmat/internal/repository"
	"github.com/xanity-07/openmat/internal/server"
)

type Service struct {
	User *UserService
	Auth *AuthService
}

func NewService(s *server.Server, repos *repository.Repositories) *Service {
	return &Service{
		User: NewUserService(s, repos.UserRepo),
		Auth: NewAuthService(s, repos.UserRepo, repos.SessionRepo),
	}
}
