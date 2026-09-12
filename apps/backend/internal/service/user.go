package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/lib/utils"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/model/user"
	"github.com/xanity-07/openmat/internal/repository"
	"github.com/xanity-07/openmat/internal/server"
)

type UserService struct {
	s        *server.Server
	UserRepo *repository.UserRepository
}

func NewUserService(s *server.Server, userRepo *repository.UserRepository) *UserService {
	return &UserService{s: s, UserRepo: userRepo}
}

func (s *UserService) CreateUser(ctx *gin.Context, payload *user.CreateUserPayload) (*user.User, error) {
	logger := middleware.GetLogger(ctx)

	exists, err := s.UserRepo.CheckUserExistsEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Str("email", payload.Email).Msg("failed to get user by email")
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if exists {
		logger.Error().Err(err).Str("email", payload.Email).Msg("user already exists")

		code := "EMAIL_ALREADY_EXISTS"
		return nil, errs.NewBadRequestError("email already in use", &code, nil, nil)
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	payload.Password = hashedPassword

	createdUser, err := s.UserRepo.CreateUser(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	logger.Info().Str("event", "user_created").
		Str("role", string(createdUser.Role)).
		Msg("User created successfully")

	return createdUser, nil
}

func (s *UserService) GetUsers(ctx *gin.Context, q *user.GetUsersQuery) ([]user.User, error) {
	logger := middleware.GetLogger(ctx)

	userList, err := s.UserRepo.GetUsers(ctx, q)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch list of users")
		return nil, err
	}

	return userList, nil
}
