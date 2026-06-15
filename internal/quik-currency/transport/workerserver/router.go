package workerserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type router struct {
	mux *chi.Mux
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func NewRouter(handler *handler) *router {
	r := chi.NewMux()
	r.Get("/health", handler.health)
	return &router{
		mux: r,
	}
}
