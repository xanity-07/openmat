package errs

import "net/http"

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewBadRequestError(message string, code *string, errors []FieldError, action *Action) *AppError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest))

	if code != nil {
		formattedCode = *code
	}

	return &AppError{
		Code:    formattedCode,
		Action:  action,
		Message: message,
		Status:  http.StatusBadRequest,
		Errors:  errors,
	}
}

func NewNotFoundError(message string, code *string) *AppError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusNotFound))

	if code != nil {
		formattedCode = *code
	}

	return &AppError{
		Code:    formattedCode,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusForbidden)),
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NewInternalServerError() *AppError {
	return &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}

func ValidateError(err error) *AppError {
	code := "VALIDATION_FAILED"
	return NewBadRequestError("Validation failed "+err.Error(), &code, nil, nil)
}
