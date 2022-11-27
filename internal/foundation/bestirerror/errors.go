package bestirerror

import (
	"fmt"
)

// WithCodeAndMessage is a convenience function for calling both
// WithStatusCode and WithUserMessage.
func WithCodeAndMessage(err error, code int, msg string) error {
	return WithStatusCode(WithUserMessage(err, msg), code)
}

// WithCodeAndMessagef is a convenience function for calling both
// WithStatusCode and WithUserMessage.
func WithCodeAndMessagef(err error, code int, format string, v ...interface{}) error {
	return WithStatusCode(WithUserMessagef(err, format, v...), code)
}

// New is a convenience function for calling fmt.Errorf and WithStatusCode.
func New(code int, format string, v ...interface{}) error {
	return WithStatusCode(
		fmt.Errorf(format, v...),
		code,
	)
}
