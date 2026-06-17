package workerserver

import (
	"context"
	"encoding/json"
	"net/http"
)

type healthDTO struct {
	Title  string `json:"title,omitempty"`
	Status int    `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type repository interface {
	PingContext(ctx context.Context) error
}

type handler struct {
	repo repository
}

func NewHandler(repo repository) *handler {
	return &handler{
		repo: repo,
	}
}

func (h handler) health(w http.ResponseWriter, r *http.Request) {

	err := h.repo.PingContext(r.Context())
	if err != nil {
		status := http.StatusServiceUnavailable
		resp := healthDTO{
			Title:  "NOT_OK",
			Status: status,
			Detail: "потеряно соединение с БД",
		}
		body, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(status)

		w.Write(body)
		return
	}
	status := http.StatusOK
	resp := healthDTO{
		Title:  "OK",
		Status: status,
		Detail: "сервис работает штатно",
	}
	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	w.Write(body)

}
