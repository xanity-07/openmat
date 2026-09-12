package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/httpResponse"
	"github.com/xanity-07/openmat/internal/model/user"
	"github.com/xanity-07/openmat/internal/server"
	"github.com/xanity-07/openmat/internal/service"
)

type UserHandler struct {
	Handler
	UserService *service.UserService
}

func NewUserHandler(s *server.Server, userService *service.UserService) *UserHandler {
	return &UserHandler{
		Handler{server: s},
		userService,
	}
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return Handle(
		h.Handler,
		func(c *gin.Context, payload *user.CreateUserPayload) (httpResponse.Response[*user.User], error) {
			createdUser, err := h.UserService.CreateUser(c, payload)
			if err != nil {
				return httpResponse.Response[*user.User]{}, err
			}

			return httpResponse.Response[*user.User]{
				Success: true,
				Status:  http.StatusCreated,
				Data:    createdUser,
			}, nil
		},
		http.StatusCreated,
		func() *user.CreateUserPayload {
			return &user.CreateUserPayload{}
		},
		enums.BindingJSON,
	)
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return Handle(
		h.Handler,
		func(c *gin.Context, q *user.GetUsersQuery) (httpResponse.Response[[]user.User], error) {
			userList, err := h.UserService.GetUsers(c, q)
			if err != nil {
				return httpResponse.Response[[]user.User]{}, err
			}

			return httpResponse.Response[[]user.User]{
				Success: true,
				Status:  http.StatusOK,
				Data:    userList,
			}, nil
		},
		http.StatusOK,
		func() *user.GetUsersQuery {
			return &user.GetUsersQuery{}
		}, enums.BindingQuery,
	)
}
