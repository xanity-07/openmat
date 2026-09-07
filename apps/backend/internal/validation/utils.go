// Package validation contains all the logic related to handler payload validations
package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/errs"
)

// Validateable whatever data type implements this Validate() error method will satisfy this interface this will be how we handle validation for all DTOs
type Validateable interface {
	Validate() error
}

type CustomValidationError struct {
	Field   string
	Message string
}

type CustomValidationErrors []CustomValidationError

// This satisfies the Error interface and now CustomValidationErrors can be of error type
func (ce CustomValidationErrors) Error() string {
	return "validation failed"
}

func BindAndValidate(c *gin.Context, payload Validateable, source enums.BindingSource) error {
	var err error

	switch source {
	case enums.BindingJSON:
		err = c.ShouldBind(payload)
	case enums.BindingQuery:
		err = c.ShouldBindQuery(payload)
	case enums.BindingURI:
		err = c.ShouldBindUri(payload)
	default:
		return errs.NewBadRequestError("invalid binding source", nil, nil, nil)
	}

	if err != nil {
		return errs.NewBadRequestError(fmt.Sprintf("binding failed: %s", err.Error()), nil, nil, nil)
	}

	if msg, fieldErrors := validateStruct(payload); fieldErrors != nil {
		code := "VALIDATION_FAILED"
		return errs.NewBadRequestError(msg, &code, fieldErrors, nil)
	}

	return nil
}

// validateStruct takes a data type that implements Validatable then we call Validate on top of that data type
// if we get any errors we call extractValidationErrors
func validateStruct(v Validateable) (string, []errs.FieldError) {
	if err := v.Validate(); err != nil {
		return extractValidationError(err)
	}
	return "", nil
}

// extractValidationErrors properly extracts the errors received from validateStruct
// so that we can send them to the client
func extractValidationError(err error) (string, []errs.FieldError) {
	var fieldErrors []errs.FieldError

	// Check if the error is from validator package
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if ok {
		// Check for validator struct tag errors
		for _, validationErr := range validationErrors {
			field := strings.ToLower(validationErr.Field())
			var msg string

			switch validationErr.Tag() {
			case "required":
				if validationErr.Type().Kind() == reflect.String {
					msg = fmt.Sprintf("%s is required", field)
				}
			case "min":
				if validationErr.Type().Kind() == reflect.String {
					msg = fmt.Sprintf("must be at least %s characters", field)
				} else {
					msg = fmt.Sprintf("must be at least %s", field)
				}
			case "max":
				if validationErr.Type().Kind() == reflect.String {
					msg = fmt.Sprintf("must not exceed %s characters", field)
				} else {
					msg = fmt.Sprintf("must be at least %s", field)
				}
			case "oneof":
				msg = fmt.Sprintf("%s must be one of %s", field, validationErr.Value())
			case "email":
				msg = fmt.Sprintf("%s must be a valid email address", field)
			case "e164":
				msg = fmt.Sprintf("%s must be a valid phone number", field)
			case "uuid":
				msg = fmt.Sprintf("%s must be a valid uuid", field)
			case "uuidList":
				msg = fmt.Sprintf("%s must be a valid comma-sepparated uuid list", field)
			case "dive":
				msg = fmt.Sprintf("%s some items are invalid", field)
			default:
				if validationErr.Param() != "" {
					msg = fmt.Sprintf("%s: %s:%s", field, validationErr.Tag(), validationErr.Param())
				} else {
					msg = fmt.Sprintf("%s: %s", field, validationErr.Tag())
				}
			}

			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: field,
				Error: msg,
			})
		}
		return "validation failed", fieldErrors
	}

	// Check if the type is CustomValidationErrors
	// loop through all the validation errors and append it into the field error slice
	var customValidationErrs CustomValidationErrors
	ok = errors.As(err, &customValidationErrs)
	if ok {
		for _, validationErr := range customValidationErrs {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: validationErr.Field,
				Error: validationErr.Message,
			})
		}
		return "validation failed", fieldErrors
	}
	return "validation failed", []errs.FieldError{
		{
			Error: err.Error(),
		},
	}
}
