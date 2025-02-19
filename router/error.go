package router

import (
	"fmt"
	"net/http"
)

var (
	_ ErrorProcessor = defaultErrorProcessor

	_ Error = (*HTTPError)(nil)
	_ Error = (*BadRequestError)(nil)
	_ Error = (*UnauthorizedError)(nil)
	_ Error = (*PaymentRequiredError)(nil)
	_ Error = (*ForbiddenError)(nil)
	_ Error = (*NotFoundError)(nil)
	_ Error = (*MethodNotAllowedError)(nil)
	_ Error = (*NotAcceptableError)(nil)
	_ Error = (*ProxyAuthRequiredError)(nil)
	_ Error = (*RequestTimeoutError)(nil)
	_ Error = (*ConflictError)(nil)
	_ Error = (*GoneError)(nil)
	_ Error = (*LengthRequiredError)(nil)
	_ Error = (*PreconditionFailedError)(nil)
	_ Error = (*RequestEntityTooLargeError)(nil)
	_ Error = (*RequestURITooLongError)(nil)
	_ Error = (*UnsupportedMediaTypeError)(nil)
	_ Error = (*RequestedRangeNotSatisfiableError)(nil)
	_ Error = (*ExpectationFailedError)(nil)
	_ Error = (*TeapotError)(nil)
	_ Error = (*MisdirectedRequestError)(nil)
	_ Error = (*UnprocessableEntityError)(nil)
	_ Error = (*LockedError)(nil)
	_ Error = (*FailedDependencyError)(nil)
	_ Error = (*TooEarlyError)(nil)
	_ Error = (*UpgradeRequiredError)(nil)
	_ Error = (*PreconditionRequiredError)(nil)
	_ Error = (*TooManyRequestsError)(nil)
	_ Error = (*RequestHeaderFieldsTooLargeError)(nil)
	_ Error = (*UnavailableForLegalReasonsError)(nil)
)

type Error interface {
	error
	StatusCode() int
}

type ErrorProcessor func(err error) error

type HTTPError struct {
	Err    error       `json:"-"`
	Title  string      `json:"title,omitzero"  description:"Short title of the error"`
	Status int         `json:"status,omitzero" description:"HTTP status code"             example:"404"`
	Detail string      `json:"detail,omitzero" description:"Human readable error message"`
	Errors []ErrorItem `json:"errors,omitzero"`
}

type ErrorItem struct {
	Name     string         `json:"name"`
	Reason   string         `json:"reason"`
	Metadata map[string]any `json:"metadata,omitzero"`
}

func (e HTTPError) Error() string {
	title := e.Title
	if title == "" {
		title = http.StatusText(e.Status)
		if title == "" {
			title = "HTTP Error"
		}
	}
	return fmt.Sprintf("%s (%d): %s", title, e.Status, e.Detail)
}

func (e HTTPError) StatusCode() int {
	if e.Status == 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e HTTPError) Unwrap() error {
	return e.Err
}

type BadRequestError HTTPError

func (e BadRequestError) Error() string   { return e.Err.Error() }
func (e BadRequestError) StatusCode() int { return http.StatusBadRequest }
func (e BadRequestError) Unwrap() error   { return HTTPError(e) }

type UnauthorizedError HTTPError

func (e UnauthorizedError) Error() string   { return e.Err.Error() }
func (e UnauthorizedError) StatusCode() int { return http.StatusUnauthorized }
func (e UnauthorizedError) Unwrap() error   { return HTTPError(e) }

type PaymentRequiredError HTTPError

func (e PaymentRequiredError) Error() string   { return e.Err.Error() }
func (e PaymentRequiredError) StatusCode() int { return http.StatusPaymentRequired }
func (e PaymentRequiredError) Unwrap() error   { return HTTPError(e) }

type ForbiddenError HTTPError

func (e ForbiddenError) Error() string   { return e.Err.Error() }
func (e ForbiddenError) StatusCode() int { return http.StatusForbidden }
func (e ForbiddenError) Unwrap() error   { return HTTPError(e) }

type NotFoundError HTTPError

func (e NotFoundError) Error() string   { return e.Err.Error() }
func (e NotFoundError) StatusCode() int { return http.StatusNotFound }
func (e NotFoundError) Unwrap() error   { return HTTPError(e) }

type MethodNotAllowedError HTTPError

func (e MethodNotAllowedError) Error() string   { return e.Err.Error() }
func (e MethodNotAllowedError) StatusCode() int { return http.StatusMethodNotAllowed }
func (e MethodNotAllowedError) Unwrap() error   { return HTTPError(e) }

type NotAcceptableError HTTPError

func (e NotAcceptableError) Error() string   { return e.Err.Error() }
func (e NotAcceptableError) StatusCode() int { return http.StatusNotAcceptable }
func (e NotAcceptableError) Unwrap() error   { return HTTPError(e) }

type ProxyAuthRequiredError HTTPError

func (e ProxyAuthRequiredError) Error() string   { return e.Err.Error() }
func (e ProxyAuthRequiredError) StatusCode() int { return http.StatusProxyAuthRequired }
func (e ProxyAuthRequiredError) Unwrap() error   { return HTTPError(e) }

type RequestTimeoutError HTTPError

func (e RequestTimeoutError) Error() string   { return e.Err.Error() }
func (e RequestTimeoutError) StatusCode() int { return http.StatusRequestTimeout }
func (e RequestTimeoutError) Unwrap() error   { return HTTPError(e) }

type ConflictError HTTPError

func (e ConflictError) Error() string   { return e.Err.Error() }
func (e ConflictError) StatusCode() int { return http.StatusConflict }
func (e ConflictError) Unwrap() error   { return HTTPError(e) }

type GoneError HTTPError

func (e GoneError) Error() string   { return e.Err.Error() }
func (e GoneError) StatusCode() int { return http.StatusGone }
func (e GoneError) Unwrap() error   { return HTTPError(e) }

type LengthRequiredError HTTPError

func (e LengthRequiredError) Error() string   { return e.Err.Error() }
func (e LengthRequiredError) StatusCode() int { return http.StatusLengthRequired }
func (e LengthRequiredError) Unwrap() error   { return HTTPError(e) }

type PreconditionFailedError HTTPError

func (e PreconditionFailedError) Error() string   { return e.Err.Error() }
func (e PreconditionFailedError) StatusCode() int { return http.StatusPreconditionFailed }
func (e PreconditionFailedError) Unwrap() error   { return HTTPError(e) }

type RequestEntityTooLargeError HTTPError

func (e RequestEntityTooLargeError) Error() string   { return e.Err.Error() }
func (e RequestEntityTooLargeError) StatusCode() int { return http.StatusRequestEntityTooLarge }
func (e RequestEntityTooLargeError) Unwrap() error   { return HTTPError(e) }

type RequestURITooLongError HTTPError

func (e RequestURITooLongError) Error() string   { return e.Err.Error() }
func (e RequestURITooLongError) StatusCode() int { return http.StatusRequestURITooLong }
func (e RequestURITooLongError) Unwrap() error   { return HTTPError(e) }

type UnsupportedMediaTypeError HTTPError

func (e UnsupportedMediaTypeError) Error() string   { return e.Err.Error() }
func (e UnsupportedMediaTypeError) StatusCode() int { return http.StatusUnsupportedMediaType }
func (e UnsupportedMediaTypeError) Unwrap() error   { return HTTPError(e) }

type RequestedRangeNotSatisfiableError HTTPError

func (e RequestedRangeNotSatisfiableError) Error() string { return e.Err.Error() }
func (e RequestedRangeNotSatisfiableError) StatusCode() int {
	return http.StatusRequestedRangeNotSatisfiable
}
func (e RequestedRangeNotSatisfiableError) Unwrap() error { return HTTPError(e) }

type ExpectationFailedError HTTPError

func (e ExpectationFailedError) Error() string   { return e.Err.Error() }
func (e ExpectationFailedError) StatusCode() int { return http.StatusExpectationFailed }
func (e ExpectationFailedError) Unwrap() error   { return HTTPError(e) }

type TeapotError HTTPError

func (e TeapotError) Error() string   { return e.Err.Error() }
func (e TeapotError) StatusCode() int { return http.StatusTeapot }
func (e TeapotError) Unwrap() error   { return HTTPError(e) }

type MisdirectedRequestError HTTPError

func (e MisdirectedRequestError) Error() string   { return e.Err.Error() }
func (e MisdirectedRequestError) StatusCode() int { return http.StatusMisdirectedRequest }
func (e MisdirectedRequestError) Unwrap() error   { return HTTPError(e) }

type UnprocessableEntityError HTTPError

func (e UnprocessableEntityError) Error() string   { return e.Err.Error() }
func (e UnprocessableEntityError) StatusCode() int { return http.StatusUnprocessableEntity }
func (e UnprocessableEntityError) Unwrap() error   { return HTTPError(e) }

type LockedError HTTPError

func (e LockedError) Error() string   { return e.Err.Error() }
func (e LockedError) StatusCode() int { return http.StatusLocked }
func (e LockedError) Unwrap() error   { return HTTPError(e) }

type FailedDependencyError HTTPError

func (e FailedDependencyError) Error() string   { return e.Err.Error() }
func (e FailedDependencyError) StatusCode() int { return http.StatusFailedDependency }
func (e FailedDependencyError) Unwrap() error   { return HTTPError(e) }

type TooEarlyError HTTPError

func (e TooEarlyError) Error() string   { return e.Err.Error() }
func (e TooEarlyError) StatusCode() int { return http.StatusTooEarly }
func (e TooEarlyError) Unwrap() error   { return HTTPError(e) }

type UpgradeRequiredError HTTPError

func (e UpgradeRequiredError) Error() string   { return e.Err.Error() }
func (e UpgradeRequiredError) StatusCode() int { return http.StatusUpgradeRequired }
func (e UpgradeRequiredError) Unwrap() error   { return HTTPError(e) }

type PreconditionRequiredError HTTPError

func (e PreconditionRequiredError) Error() string   { return e.Err.Error() }
func (e PreconditionRequiredError) StatusCode() int { return http.StatusPreconditionRequired }
func (e PreconditionRequiredError) Unwrap() error   { return HTTPError(e) }

type TooManyRequestsError HTTPError

func (e TooManyRequestsError) Error() string   { return e.Err.Error() }
func (e TooManyRequestsError) StatusCode() int { return http.StatusTooManyRequests }
func (e TooManyRequestsError) Unwrap() error   { return HTTPError(e) }

type RequestHeaderFieldsTooLargeError HTTPError

func (e RequestHeaderFieldsTooLargeError) Error() string { return e.Err.Error() }
func (e RequestHeaderFieldsTooLargeError) StatusCode() int {
	return http.StatusRequestHeaderFieldsTooLarge
}
func (e RequestHeaderFieldsTooLargeError) Unwrap() error { return HTTPError(e) }

type UnavailableForLegalReasonsError HTTPError

func (e UnavailableForLegalReasonsError) Error() string { return e.Err.Error() }
func (e UnavailableForLegalReasonsError) StatusCode() int {
	return http.StatusUnavailableForLegalReasons
}
func (e UnavailableForLegalReasonsError) Unwrap() error { return HTTPError(e) }

func defaultErrorProcessor(err error) error {
	errResponse := HTTPError{
		Err:    err,
		Status: http.StatusInternalServerError,
		Detail: "An unexpected error occurred",
	}

	if httpErr, ok := err.(HTTPError); ok {
		errResponse = httpErr
	}

	if errorStatus, ok := err.(Error); ok {
		errResponse.Status = errorStatus.StatusCode()
	}

	if errResponse.Title == "" {
		errResponse.Title = http.StatusText(errResponse.Status)
	}

	return errResponse
}
