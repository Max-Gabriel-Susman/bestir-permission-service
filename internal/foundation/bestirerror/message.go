package bestirerror

import (
	"errors"
	"fmt"
	"net/http"
)

// UserMessenger is an error with an associated user-facing message.
type UserMessenger interface {
	error
	UserMessage() string
}

type messenger struct {
	error
	msg string
}

func (msgr messenger) Unwrap() error {
	return msgr.error
}

func (msgr messenger) UserMessage() string {
	return msgr.msg
}

// WithUserMessage adds a UserMessenger to err's error chain.
// This func will wrap a nil error.
func WithUserMessage(err error, msg string) error {
	if err == nil {
		err = errors.New(msg)
	}
	return messenger{err, msg}
}

// WithUserMessagef calls fmt.Sprintf before calling WithUserMessage.
func WithUserMessagef(err error, format string, v ...interface{}) error {
	return WithUserMessage(err, fmt.Sprintf(format, v...))
}

// UserMessage returns the user message associated with an error.
// If no message is found, it checks StatusCode and returns that message.
// Because the default status is 500, the default message is
// "Internal Server Error".
// If err is nil, it returns "".
func UserMessage(err error) string {
	if err == nil {
		return ""
	}
	if um := UserMessenger(nil); errors.As(err, &um) {
		return um.UserMessage()
	}
	return http.StatusText(StatusCode(err))
}
