package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/server"
)

type OpenAPIHandler struct {
	Handler
}

func NewOpenAPIHandler(s *server.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		Handler{server: s},
	}
}

func (h *OpenAPIHandler) ServerOpenAPIUI(c *gin.Context) {
	templateBytes, err := os.ReadFile("static/openapi.html")
	if err != nil {
		h.server.Logger.Error().Msgf("Failed to server openapi.html: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Do not cache the file so we can see new changes with refresh
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store")
	c.Data(http.StatusOK, "text/html; charset=utf-8", templateBytes)
}
