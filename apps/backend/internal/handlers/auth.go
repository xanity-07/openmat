// internal/handlers/auth.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/httpResponse"
	"github.com/xanity-07/openmat/internal/model/authmodel"
	"github.com/xanity-07/openmat/internal/server"
	"github.com/xanity-07/openmat/internal/service"
)

type AuthHandler struct {
	Handler
	AuthService *service.AuthService
}

func NewAuthHandler(s *server.Server, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{Handler{server: s}, authService}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return Handle(
		h.Handler,
		func(c *gin.Context, payload *authmodel.LoginRequest) (httpResponse.Response[string], error) {
			token, err := h.AuthService.Login(c, payload)
			if err != nil {
				return httpResponse.Response[string]{}, err
			}

			return httpResponse.Response[string]{
				Success: true,
				Status:  http.StatusOK,
				Data:    token,
			}, nil
		},
		http.StatusOK,
		func() *authmodel.LoginRequest { return &authmodel.LoginRequest{} },
		enums.BindingJSON,
	)
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return HandleNoContent(
		h.Handler,
		func(c *gin.Context, payload EmptyRequest) error {
			return h.AuthService.Logout(c)
		},
		http.StatusNoContent,
		func() EmptyRequest { return EmptyRequest{} },
		enums.BindingJSON,
	)
}
