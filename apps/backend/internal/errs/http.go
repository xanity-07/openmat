// Package errs contains all the logic related to application errors and API responses
package errs

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ActionType string

const (
	ActionTypeRedirect ActionType = "redirect"
)

// Action handles what action we take in specific API requests like redirects etc
type Action struct {
	Type    ActionType `json:"type"`
	Message string     `json:"message"`
	Value   string     `json:"value"`
}

type AppError struct {
	Action  *Action      `json:"action"`
	Errors  []FieldError `json:"errors"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Status  int          `json:"status"`
}

func (e *AppError) Error() string {
	return e.Message
}

type Response[T any] struct {
	Success bool       `json:"success"`
	Status  int        `json:"status"`
	Data    T          `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorInfo struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors,omitempty"`
	Action  *Action      `json:"action,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

// MakeUpperCaseWithUnderscores returns a HTTP status with format example: "BAD_REQUEST"
func MakeUpperCaseWithUnderscores(status string) string {
	return strings.ToUpper(strings.ReplaceAll(status, " ", "_"))
}

func WriteHTTPError(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		response := Response[any]{
			Success: false,
			Status:  appErr.Status,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Errors:  appErr.Errors,
				Action:  appErr.Action,
			},
		}

		c.AbortWithStatusJSON(appErr.Status, response)
		return
	}

	response := Response[any]{
		Success: false,
		Status:  http.StatusInternalServerError,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		},
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, response)
}
