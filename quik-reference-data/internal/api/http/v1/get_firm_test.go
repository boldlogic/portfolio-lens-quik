package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type readSvcStub struct{}

func (readSvcStub) GetFirms(ctx context.Context) ([]quik.Firm, error) {
	return nil, nil
}

func (readSvcStub) GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	return quik.Firm{}, nil
}

type writeSvcStub struct{}

func (writeSvcStub) CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	return quik.Firm{}, nil
}

func (writeSvcStub) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	return quik.Firm{}, nil
}

func reqFirmWithID(id string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/firms/"+id, nil)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rc))
}

func TestGetFirm_некорректный_id_ErrValidation(t *testing.T) {
	h := NewHandler(handler.NewHandler(), readSvcStub{}, writeSvcStub{}, zap.NewNop())
	_, _, err := h.GetFirm(reqFirmWithID("x"))
	if !errors.Is(err, models.ErrValidation) {
		t.Fatalf("ожидали ErrValidation: %v", err)
	}
}
