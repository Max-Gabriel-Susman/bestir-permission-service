package bestirerror_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/bestirerror"
)

func describe(err error) {
	fmt.Printf("String: %s\nUserMessage: %s\nStatusCode: %d\nDetails: %+v",
		err, bestirerror.UserMessage(err), bestirerror.StatusCode(err), bestirerror.Details(err))
}

func ExampleStatusCoder() {
	err := errors.New("example 404 error")
	errCode := bestirerror.WithStatusCode(err, http.StatusNotFound)
	describe(errCode)
	// Output:
	// String: [404] example 404 error
	// UserMessage: Not Found
	// StatusCode: 404
	// Details: []
}

func ExampleUserMessenger() {
	err := errors.New("example 404 error")
	errCode := bestirerror.WithUserMessage(err, "This is the message returned to the user")
	describe(errCode)
	// Output:
	// String: example 404 error
	// UserMessage: This is the message returned to the user
	// StatusCode: 500
	// Details: []
}

func ExampleWithCodeAndMessage() {
	err := errors.New("example 404 error")
	errCodeMsg := bestirerror.WithCodeAndMessage(err, http.StatusNotFound, "This is the message returned to the user")
	describe(errCodeMsg)
	// Output:
	// String: [404] example 404 error
	// UserMessage: This is the message returned to the user
	// StatusCode: 404
	// Details: []
}

func ExampleDetailer() {
	err := errors.New("example 404 error")
	errCodeMsg := bestirerror.WithDetails(err, []string{"I am some details", "yay"})
	describe(errCodeMsg)
	// Output:
	// String: example 404 error
	// UserMessage: Internal Server Error
	// StatusCode: 500
	// Details: [I am some details yay]
}
