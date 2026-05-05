package v1

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/go-chi/chi/v5"
)

type readSvc struct {
	err error
}

func (s readSvc) GetFirms(ctx context.Context) ([]quik.Firm, error) {
	return nil, nil
}

func (s readSvc) GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	return quik.Firm{}, nil
}

func reqFirmWithID(method string, id string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, "/firms/"+id, body)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rc))
}
