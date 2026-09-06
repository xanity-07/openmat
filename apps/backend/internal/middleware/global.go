// Package middleware contains all the middleware functions in our application
package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/server"
	"github.com/xanity-07/openmat/internal/sqlerr"
)

type GlobalMiddlewares struct {
	server *server.Server
}

func NewGlobalMiddleware(s *server.Server) *GlobalMiddlewares {
	return &GlobalMiddlewares{
		server: s,
	}
}

func (global *GlobalMiddlewares) CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: global.server.Config.Server.CORSAllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	})
}

func (global *GlobalMiddlewares) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Execute the rest of the middlewares/downstream layers
		c.Next()

		latency := time.Since(start)

		statusCode := c.Writer.Status()
		ginErr := c.Errors.Last()

		if ginErr != nil {
			if appErr, ok := errors.AsType[*errs.AppError](ginErr.Err); ok {
				statusCode = appErr.Status
			}
		}

		logger := GetLogger(c)

		var event *zerolog.Event

		switch {
		case statusCode >= 500:
			event = logger.Error()
			if ginErr != nil {
				event = event.Err(ginErr.Err)
			}

		case statusCode >= 400:
			event = logger.Error()
			if ginErr != nil {
				event = event.Err(ginErr.Err)
			}

		default:
			event = logger.Info()
		}

		// Add request ID from context if available
		if requestID := GetRequestID(c); requestID != "" {
			event.Str("request_id", requestID)
		}

		// Add user ID from context if available
		if userID := GetUserID(c); userID != "" {
			event.Str("user_id", userID)
		}

		// Add user role from context if available
		if userRole := GetUserRole(c); userRole != "" {
			event.Str("user_role", userRole)
		}

		// Add session ID from context if available
		if sessionID := GetSessionID(c); sessionID != "" {
			event.Str("session_id", sessionID)
		}

		event.
			Dur("duration", latency).
			Int("status", statusCode).
			Str("method", c.Request.Method).
			Str("uri", c.Request.RequestURI).
			Str("host", c.Request.Host).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("API")
	}
}

func (global *GlobalMiddlewares) Recover() gin.HandlerFunc {
	return gin.CustomRecovery(
		func(c *gin.Context, err any) {
			logger := GetLogger(c)

			logger.Error().
				Any("panic", err).
				Msg("panic recovered")

			c.AbortWithStatusJSON(http.StatusInternalServerError, errs.Response[any]{
				Success: false,
				Status:  http.StatusInternalServerError,
				Data:    nil,
				Error: &errs.ErrorInfo{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "internal server error",
				},
			})
		},
	)
}

func (global *GlobalMiddlewares) Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")                                          // Click jacking via iframes
		c.Header("X-Content-Type-Options", "nosniff")                                // MIME-type sniffing attacks
		c.Header("Content-Security-Policy", "default-src 'self'")                    // XSS and data injection
		c.Header("Referrer-Policy", "strict-origin")                                 // Leaking sensitive URL parameters to third parties
		c.Header("Permissions-Policy", "geolocation=(), camera=(), microphone=()")   // Unauthorized use of browser APIs (camera, mic, etc.)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains") // Protocol downgrade and cookie hijacking
		c.Next()
	}
}

func (global *GlobalMiddlewares) GlobalErrorHandling(c *gin.Context, err error) {
	// First try to handle database errors and convert them into appropriate HTTP errors
	originalError := err

	// Handling well known database errors
	// Only for errors that haven't been converted to a application error
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		err = errs.NewBadRequestError(err.Error(), nil, nil, nil)
	} else {
		// Here we call our sqlerr handler that will convert database errors to appropriate application errors
		err = sqlerr.HandleError(err)
	}

	// Now process the possible converted error
	var status int
	var code string
	var message string
	var fieldError []errs.FieldError
	var action *errs.Action

	switch {
	case errors.As(err, &appErr):
		status = appErr.Status
		code = appErr.Code
		message = appErr.Message
		fieldError = append(fieldError, appErr.Errors...)
		action = appErr.Action

	default:
		status = http.StatusInternalServerError
		code = errs.MakeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError))
		message = http.StatusText(http.StatusInternalServerError)
	}

	// Log the original error to help with debugging
	// Use enhanced logger from context which already includes request_id, method, path, ip, user_role, user_id, session_id, and trace context
	logger := *GetLogger(c)

	logger.Error().Stack().
		Err(originalError).
		Int("status", status).
		Str("error_code", code).
		Msg(message)

	if !c.Writer.Written() {
		c.JSON(status, errs.Response[any]{
			Success: false,
			Status:  status,
			Error: &errs.ErrorInfo{
				Code:    code,
				Message: message,
				Errors:  fieldError,
				Action:  action,
			},
		})
	}
}
