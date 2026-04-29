package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type svcStub struct{}

func (svcStub) GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	return nil, nil
}

func (svcStub) GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	return nil, nil
}

func (svcStub) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	return nil, nil
}

func (svcStub) CreateMoneyLimit(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {
	return quik.MoneyLimit{}, nil
}

func (svcStub) CreateSecurityLimit(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
	return quik.SecurityLimit{}, nil
}

func (svcStub) CreateSecurityLimitOtc(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
	return quik.SecurityLimit{}, nil
}

func (svcStub) GetLimits(ctx context.Context, date time.Time) ([]quik.Limit, error) {
	return nil, nil
}

func (svcStub) GetPortfolio(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error) {
	return nil, nil
}

func (svcStub) GetFirms(ctx context.Context) ([]quik.Firm, error) {
	return nil, nil
}

func (svcStub) GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	return quik.Firm{}, nil
}

func (svcStub) CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	return quik.Firm{}, nil
}

func (svcStub) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	return quik.Firm{}, nil
}

func reqFirmWithID(id string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/firms/"+id, nil)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rc))
}

func TestGetFirm_некорректный_id_ErrValidation(t *testing.T) {
	h := NewHandler(handler.NewHandler(), svcStub{}, zap.NewNop())
	_, _, err := h.GetFirm(reqFirmWithID("x"))
	if !errors.Is(err, models.ErrValidation) {
		t.Fatalf("ожидали ErrValidation: %v", err)
	}
}
