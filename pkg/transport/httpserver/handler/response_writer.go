package handler

import (
	"encoding/json"
	"net/http"
)

type HTTPErr struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func WriteResp(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func NotFound(detail string) HTTPErr {
	return HTTPErr{
		Title:  "NOT_FOUND",
		Status: http.StatusNotFound,
		Detail: detail,
	}
}

func BadRequest(detail string) HTTPErr {
	return HTTPErr{
		Title:  "VALIDATION_ERROR",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func UnprocessableEntity(detail string) HTTPErr {
	return HTTPErr{
		Title:  "BUSINESS_VALIDATION_ERROR",
		Status: http.StatusUnprocessableEntity,
		Detail: detail,
	}
}
func Internal(detail string) HTTPErr {
	return HTTPErr{
		Title:  "SERVER_ERROR",
		Status: http.StatusInternalServerError,
		Detail: "что-то пошло не так",
	}
}

func Conflict(detail string) HTTPErr {
	return HTTPErr{
		Title:  "CONFLICT",
		Status: http.StatusConflict,
		Detail: detail,
	}
}

func UnsupportedMediaType(detail string) HTTPErr {
	return HTTPErr{
		Title:  "UNSUPPORTED_MEDIA_TYPE",
		Status: http.StatusUnsupportedMediaType,
		Detail: detail,
	}
}

func RequestEntityTooLarge(detail string) HTTPErr {
	return HTTPErr{
		Title:  "REQUEST_ENTITY_TOO_LARGE",
		Status: http.StatusRequestEntityTooLarge,
		Detail: detail,
	}
}
