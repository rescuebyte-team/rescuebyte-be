package errx

import (
	"github.com/gofiber/fiber/v3"
)

type RequestError struct {
	StatusCode int            `json:"-"`
	ErrorCode  string         `json:"error_code"`
	Message    string         `json:"message"`
	Location   string         `json:"location,omitempty"`
	Details    map[string]any `json:"details,omitempty"`
	Err        error          `json:"-"`
}

func (r *RequestError) Error() string {
	return r.Message
}

func (r *RequestError) Is(target error) bool {
	t, ok := target.(*RequestError)
	if !ok {
		return false
	}
	return r.ErrorCode == t.ErrorCode
}

func NewError(statusCode int, errorCode string, message string) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (r *RequestError) clone() *RequestError {
	c := *r
	return &c
}

func (r *RequestError) WithLocation(location string) *RequestError {
	c := r.clone()
	c.Location = location
	return c
}

func (r *RequestError) WithMessage(message string) *RequestError {
	c := r.clone()
	c.Message = message
	return c
}

func (r *RequestError) WithDetails(details map[string]any) *RequestError {
	c := r.clone()
	c.Details = details
	return c
}

func (r *RequestError) WithError(err error) *RequestError {
	c := r.clone()
	c.Err = err
	return c
}

var (
	ErrInternalServer = NewError(
		fiber.StatusInternalServerError,
		"internal_server_error",
		"Something went wrong on our end. Please try again later.",
	)
	ErrNotFound = NewError(
		fiber.StatusNotFound,
		"not_found",
		"The requested resource could not be found.",
	)
	ErrForbidden = NewError(
		fiber.StatusForbidden,
		"forbidden",
		"You do not have permission to access this resource.",
	)
	ErrUnauthorized = NewError(
		fiber.StatusUnauthorized,
		"unauthorized",
		"You are not authorized to access this resource.",
	)
	ErrBadRequest = NewError(
		fiber.StatusBadRequest,
		"bad_request",
		"The request payload or parameter is not valid.",
	)
	ErrNoBearerToken = NewError(
		fiber.StatusUnauthorized,
		"no_bearer_token",
		"A bearer token is required to continue.",
	)
	ErrInvalidBearerToken = NewError(
		fiber.StatusUnauthorized,
		"invalid_bearer_token",
		"The bearer token provided is invalid or expired.",
	)
	ErrClaimsNotFound = NewError(
		fiber.StatusUnauthorized,
		"claims_not_found",
		"Authorization details are missing. Please try again.",
	)
	ErrInvalidTimeFormat = NewError(
		fiber.StatusBadRequest,
		"invalid_time_format",
		"Invalid time format. Please use HH:mm or HH.mm format (e.g. 08:00 or 08.00).",
	)
)
