package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const MaxRequestBodySize = 64 * 1024

var (
	ErrWrongJSON             = errors.New("некорректный формат запроса")
	ErrUnsupportedMediaType  = errors.New("Content-Type должен быть application/json")
	ErrRequestEntityTooLarge = errors.New("тело запроса превышает ограничение")
)

func checkContentType(r *http.Request) error {
	ct := strings.TrimSpace(r.Header.Get("Content-Type"))
	if ct == "" {
		return ErrUnsupportedMediaType
	}
	if !strings.HasPrefix(strings.ToLower(ct), "application/json") {
		return ErrUnsupportedMediaType
	}
	return nil
}

func DecodeJSON[T any](r *http.Request) (T, error) {
	var v T

	if err := checkContentType(r); err != nil {
		return v, err
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, MaxRequestBodySize+1))
	if err != nil {
		return v, ErrWrongJSON
	}
	if len(body) > MaxRequestBodySize {
		return v, ErrRequestEntityTooLarge
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		return v, ErrWrongJSON
	}
	if decoder.More() {
		return v, ErrWrongJSON
	}
	return v, nil
}

func DecodeAndValidate[T any](r *http.Request) (T, error) {
	v, err := DecodeJSON[T](r)
	if err != nil {
		return v, err
	}
	if err := validate.Struct(v); err != nil {
		return v, err
	}
	return v, nil
}

var validate *validator.Validate

func init() {

	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateJSONStruct(req any) error {
	return validate.Struct(req)
}
