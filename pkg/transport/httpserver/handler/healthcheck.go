package handler

import "net/http"

func (h *Handler) Healthcheck(r *http.Request) (any, string, error) {
	return "ok", "", nil
}
