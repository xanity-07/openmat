package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/middleware"
	"github.com/xanity-07/openmat/internal/server"
	"github.com/xanity-07/openmat/internal/validation"
)

// Handler provides base functionality for all handlers
type Handler struct {
	server *server.Server
}

// NewHandler creates a new base handler
func NewHandler(s *server.Server) Handler {
	return Handler{server: s}
}

// EmptyRequest will satisfy the compiler when implementing user logout
type EmptyRequest struct{}

func (EmptyRequest) Validate() error { return nil }

// HandleFunc represents a typed handler function that processes a request and returns a response
type HandleFunc[Req validation.Validateable, Res any] func(c *gin.Context, payload Req) (Res, error)

// HandleFuncNoContent represents a types handler function that processes a request and returns no content
type HandleFuncNoContent[Req validation.Validateable] func(c *gin.Context, payload Req) error

// ResponseHandler defines the interface for handling different response types
type ResponseHandler interface {
	Handle(c *gin.Context, result any)
	GetOperation() string
	AddAttribute(txn *newrelic.Transaction, result any)
}

// JSONResponseHandler implements ResponseHandler and handles only JSON responses
type JSONResponseHandler struct {
	status int
}

func (h JSONResponseHandler) Handle(c *gin.Context, result any) {
	c.JSON(h.status, result)
}

func (h JSONResponseHandler) GetOperation() string {
	return "handler"
}

func (h JSONResponseHandler) AddAttribute(txn *newrelic.Transaction, result any) {
	// http.status_code is already set by tracing middleware
}

// NoContentResponseHandler handles no-content responses
type NoContentResponseHandler struct {
	status int
}

func (h NoContentResponseHandler) Handle(c *gin.Context, result any) {
	c.Status(http.StatusNoContent)
}

func (h NoContentResponseHandler) GetOperation() string {
	return "handler"
}

func (h NoContentResponseHandler) AddAttribute(txn *newrelic.Transaction, result any) {
	// http.status_code is already set by tracing middleware
}

// handleRequest is the unified handler function that eliminates code duplication
func handleRequest[Req validation.Validateable](
	c *gin.Context,
	req Req,
	handler func(c *gin.Context, req Req) (interface{}, error),
	responseHandler ResponseHandler,
	source enums.BindingSource,
) error {
	start := time.Now()
	method := c.Request.Method
	path := c.Request.URL.Path
	route := path

	// Get New Relic transaction from context
	txn := newrelic.FromContext(c.Request.Context())
	if txn != nil {
		txn.AddAttribute("handler.name", route)
		responseHandler.AddAttribute(txn, nil)
	}

	// Get context-enhancer from context
	logger := middleware.GetLogger(c).With().
		Str("operation", responseHandler.GetOperation()).
		Str("method", method).
		Str("path", path).
		Str("route", route).
		Logger()

	logger.Info().Msg("Handling request")

	// Validation with observability
	validationStart := time.Now()
	if err := validation.BindAndValidate(c, req, source); err != nil {
		validationDuration := time.Since(validationStart)

		logger.Error().Err(err).
			Dur("validation_duration", validationDuration).
			Msg("Validation failed")

		if txn != nil {
			txn.NoticeError(nrpkgerrors.Wrap(err))
			txn.AddAttribute("validation.status", "failed")
			txn.AddAttribute("validation.duration_ms", validationDuration.Milliseconds())
		}
		errs.WriteHTTPError(c, err)
		return err
	}

	validationDuration := time.Since(validationStart)
	if txn != nil {
		txn.AddAttribute("validation.status", "succeeded")
		txn.AddAttribute("validation.duration_ms", validationDuration.Milliseconds())
	}

	logger.Info().
		Dur("validation_duration", validationDuration).
		Msg("request validation success")

	// Execute handler with observability
	handlerStart := time.Now()
	result, err := handler(c, req)
	handlerDuration := time.Since(handlerStart)
	if err != nil {
		totalDuration := time.Since(start)

		logger.Error().Err(err).
			Dur("handler_duration", time.Duration(handlerDuration.Milliseconds())).
			Dur("total_duration", time.Duration(totalDuration.Milliseconds())).
			Msg("handler execution failed")

		if txn != nil {
			txn.NoticeError(nrpkgerrors.Wrap(err))
			txn.AddAttribute("handler.status", "failed")
			txn.AddAttribute("handler.duration_ms", handlerDuration.Milliseconds())
			txn.AddAttribute("total.duration_ms", totalDuration.Milliseconds())
		}

		errs.WriteHTTPError(c, err)
		return err
	}

	totalDuration := time.Since(start)

	// Record success metrics and tracing
	if txn != nil {
		txn.AddAttribute("handler.status", "succeeded")
		txn.AddAttribute("handler.duration_ms", handlerDuration.Milliseconds())
		txn.AddAttribute("total.duration_ms", totalDuration.Milliseconds())
		responseHandler.AddAttribute(txn, result)
	}

	logger.Info().
		Dur("handler_duration", handlerDuration).
		Dur("validation_duration", validationDuration).
		Dur("total_duration", totalDuration).
		Msg("Handler execution success")

	responseHandler.Handle(c, result)

	return nil
}

// Handle wraps a handler with validation, error handler, logging, metrics and tracing
func Handle[Req validation.Validateable, Res any](
	h Handler,
	handler HandleFunc[Req, Res],
	status int,
	payload func() Req,
	source enums.BindingSource,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handleRequest(
			c,
			payload(),
			func(c *gin.Context, req Req) (interface{}, error) {
				return handler(c, req)
			},
			JSONResponseHandler{status: status},
			source,
		)
	}
}

// HandleNoContent wraps a handler with validation, error handling, logging, metrics, and tracing for endpoints that return no content
func HandleNoContent[Req validation.Validateable](
	h Handler,
	handler HandleFuncNoContent[Req],
	status int,
	payload func() Req,
	source enums.BindingSource,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handleRequest(
			c,
			payload(),
			func(c *gin.Context, req Req) (interface{}, error) {
				return nil, handler(c, req)
			},
			NoContentResponseHandler{status: status},
			source,
		)
	}
}
