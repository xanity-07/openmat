package httpResponse

import (
	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/model"
)

type Response[T any] struct {
	Success bool            `json:"success"`
	Status  int             `json:"status"`
	Data    T               `json:"data,omitempty"`
	Error   *errs.ErrorInfo `json:"error,omitempty"`
	Meta    *model.Meta     `json:"meta,omitempty"`
}
