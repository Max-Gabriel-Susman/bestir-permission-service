package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/ettle/strcase"
	"github.com/go-playground/validator/v10"

	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/bestirerror"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func Decode(r io.Reader, val interface{}) error {
	err := json.NewDecoder(r).Decode(val)
	if err == nil {
		return Validate(val)
	}
	if errors.Is(err, io.EOF) {
		return bestirerror.WithCodeAndMessage(err, http.StatusBadRequest, "required request body was not provided")
	}

	return bestirerror.WithCodeAndMessage(err, http.StatusBadRequest, err.Error())
}

func Validate(val interface{}) error {
	once.Do(func() {
		validate = validator.New()
	})

	if err := validate.Struct(val); err != nil {
		var validationErrors validator.ValidationErrors
		if !errors.As(err, &validationErrors) {
			return err
		}
		det := make([]string, 0, len(validationErrors))
		for i := range validationErrors {
			v := validationErrors[i]
			param := v.Param()
			if param != "" {
				param = " " + param
			}

			det = append(det, strcase.ToSnake(v.Field())+" is required"+param)
		}
		return ErrValidation(err, det)
	}
	return nil
}

func ErrValidation(e error, details []string) error {
	err := bestirerror.WithCodeAndMessage(e, http.StatusBadRequest, "request validation error")
	return bestirerror.WithDetails(err, details)
}
