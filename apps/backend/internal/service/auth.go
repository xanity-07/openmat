// internal/service/auth.go
package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/lib/utils"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/model/authmodel"
	"github.com/xanity-07/openmat/internal/model/session"
	"github.com/xanity-07/openmat/internal/repository"
	"github.com/xanity-07/openmat/internal/server"
)

type AuthService struct {
	s           *server.Server
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

func NewAuthService(s *server.Server, userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{s: s, UserRepo: userRepo, SessionRepo: sessionRepo}
}

func (s *AuthService) Login(ctx *gin.Context, payload *authmodel.LoginRequest) (string, error) {
	logger := middleware.GetLogger(ctx)

	dbUser, err := s.UserRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Str("email", payload.Email).Msg("login failed: user lookup")
		// Deliberately vague — don't leak "email not found" vs "wrong password"
		return "", errs.NewUnauthorizedError("invalid email or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(payload.Password)); err != nil {
		return "", errs.NewUnauthorizedError("invalid email or password")
	}

	sessionID := utils.GenerateID(32)
	ttl := time.Duration(s.s.Config.Auth.TTLHours) * time.Hour

	if err = s.SessionRepo.Create(ctx, sessionID, session.Session{
		UserID: dbUser.ID,
		Role:   dbUser.Role,
	}, ttl); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	token, err := utils.GenerateJWT(s.s.Config.Auth.SecretKey, dbUser.ID, sessionID, string(dbUser.Role), ttl)
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}
func (s *AuthService) Logout(ctx *gin.Context) error {
	sessionID := middleware.GetSessionID(ctx)
	if sessionID == "" {
		return errs.NewUnauthorizedError("not authenticated")
	}

	if err := s.SessionRepo.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
