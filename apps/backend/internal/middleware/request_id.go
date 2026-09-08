package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/lib/utils"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// RequestID identifier to every request so the whole request life-cycle so interactions from all our components can be grouped into one transaction
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		if requestID == "" {
			requestID = utils.GenerateID(11)
		}

		// Set the requestID in our context for down stream layers
		c.Set(RequestIDKey, requestID)
		// Set X-Request-ID header to our newly generated requestID
		c.Writer.Header().Set(RequestIDHeader, requestID)
	}
}

// GetRequestID is a helper function that retrieves the current requests ID from our context
func GetRequestID(c *gin.Context) string {
	if requestID, ok := c.Get(RequestIDKey); ok {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
