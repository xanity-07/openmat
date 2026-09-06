package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/xanity-07/openmat/internal/loggerpkg"
	"github.com/xanity-07/openmat/internal/server"
)

type CTXKey string

const (
	UserIDKey    CTXKey = "user_id"
	UserRoleKey  CTXKey = "user_role"
	LoggerKey    CTXKey = "logger"
	SessionIDKey CTXKey = "session_id"
)

type ContextEnhancer struct {
	server *server.Server
}

func NewContextEnhancer(s *server.Server) *ContextEnhancer {
	return &ContextEnhancer{
		server: s,
	}
}

func (ce *ContextEnhancer) EnhanceContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract request ID
		requestID := GetRequestID(c)

		// Create an enhance logger with our request EnhanceContext
		contextLogger := ce.server.Logger.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Logger()

		// Add trace context if available
		if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
			contextLogger = loggerpkg.WithTraceContext(contextLogger, txn)
		}

		// Extract user information from JWT or Session
		if userID := ce.extractUserID(c); userID != "" {
			contextLogger.With().
				Str("user_id", userID).
				Logger()
		}

		if userRole := ce.extractUserRole(c); userRole != "" {
			contextLogger.With().
				Str("user_role", userRole).
				Logger()
		}

		// Store the enhanced logger in the context so that in our downstream layers
		// handler -> service -> repository etc, can use the same instance of the logger
		// so that all the components it touched are grouped to be part of the same interaction
		c.Set(LoggerKey, &contextLogger)

		// Create a new context and store the logger inside it
		ctx := context.WithValue(c.Request.Context(), LoggerKey, &contextLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// extractUserID checks if user_role is set by auth middleware and extracts it
func (ce *ContextEnhancer) extractUserID(c *gin.Context) string {
	if userID, ok := c.Get(UserIDKey); ok {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// extractUserRole checks if user_role is set by auth middleware and extracts it
func (ce *ContextEnhancer) extractUserRole(c *gin.Context) string {
	if userRole, ok := c.Get(UserRoleKey); ok {
		if role, ok := userRole.(string); ok {
			return role
		}
	}
	return ""
}

// GetUserID is a utility function that returns the userID from our context
func GetUserID(c *gin.Context) string {
	if userID, ok := c.Get(UserIDKey); ok {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserRole is a utility function that returns the userRole from our context
func GetUserRole(c *gin.Context) string {
	if userRole, ok := c.Get(UserRoleKey); ok {
		if role, ok := userRole.(string); ok {
			return role
		}
	}
	return ""
}

// GetLogger is a utility function that returns the contextLogger from our context
func GetLogger(c *gin.Context) *zerolog.Logger {
	if loggerID, ok := c.Get(LoggerKey); ok {
		if logger, ok := loggerID.(*zerolog.Logger); ok {
			return logger
		}
	}

	logger := zerolog.Nop()
	return &logger
}

// GetSessionID is a utility function that returns the session ID from our context
func GetSessionID(c *gin.Context) string {
	if sessionID, ok := c.Get(SessionIDKey); ok {
		if sID, ok := sessionID.(string); ok {
			return sID
		}
	}
	return ""
}
