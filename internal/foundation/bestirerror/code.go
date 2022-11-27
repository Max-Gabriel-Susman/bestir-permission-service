package bestirerror

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// StatusCoder is an error with an associated HTTP status code.
type StatusCoder interface {
	error
	StatusCode() int
}

type statusCoder struct {
	error
	code int
}

func (sc statusCoder) Unwrap() error {
	return sc.error
}

func (sc statusCoder) Error() string {
	return fmt.Sprintf("[%d] %v", sc.code, sc.error)
}

func (sc statusCoder) StatusCode() int {
	return sc.code
}

// WithStatusCode adds a StatusCoder to err's error chain.
// (This function will wrap nil errors)
// This function also keeps you from overriding the code for any of our "known" error types.
func WithStatusCode(err error, code int) error {
	if err == nil {
		err = errors.New(http.StatusText(code))
	}
	code = checkShouldOverride(err, code)
	return statusCoder{err, code}
}

// StatusCode returns the status code associated with an error.
// If no status code is found, it returns 500 http.StatusInternalServerError.
// As a special case, it checks for Timeout() and Temporary errors and returns
// 504 http.StatusGatewayTimeout and 503 http.StatusServiceUnavailable respectively.
// If the error is a context Canceled error, it returns status code 499 (stolen from nginx).
// If err is nil, it returns 200 http.StatusOK.
func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if sc := StatusCoder(nil); errors.As(err, &sc) {
		return sc.StatusCode()
	}
	return checkShouldOverride(err, http.StatusInternalServerError)
}

// checkShouldOverride takes an error and a "default" code.
// It checks the error to see if is of a certain known type (context Canceled, not found etc...)
// and if so it overides the given status code with the appropriate code to permission the type.
func checkShouldOverride(err error, code int) int {
	var timeouter interface {
		error
		Timeout() bool
	}
	if errors.As(err, &timeouter) && timeouter.Timeout() {
		return http.StatusGatewayTimeout
	}
	var temper interface {
		error
		Temporary() bool
	}
	if errors.As(err, &temper) && temper.Temporary() {
		return http.StatusServiceUnavailable
	}
	if errors.Is(err, context.Canceled) {
		// HTTP 499 in Nginx means that the client closed the connection before the server answered the request.
		// usually caused by client side timeout.
		return 499
	}
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound
	}
	return code
}
