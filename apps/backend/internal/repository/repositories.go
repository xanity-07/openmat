package repository

import "github.com/xanity-07/openmat/internal/server"

type Repositories struct {
	UserRepo    *UserRepository
	SessionRepo *SessionRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		UserRepo:    NewUserRepository(s),
		SessionRepo: NewSessionRepository(s),
	}
}
