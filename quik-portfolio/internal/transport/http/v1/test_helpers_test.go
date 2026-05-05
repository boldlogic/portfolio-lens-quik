package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type svc struct {
	err                    error
	getLimits              func(ctx context.Context, date time.Time) ([]quik.Limit, error)
	getMoneyLimits         func(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error)
	getSecurityLimits      func(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	getSecurityLimitsOtc   func(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	getPortfolio           func(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error)
	createMoneyLimit       func(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error)
	createSecurityLimit    func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error)
	createSecurityLimitOtc func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error)
}

func (s svc) GetLimits(ctx context.Context, date time.Time) ([]quik.Limit, error) {
	if s.getLimits != nil {
		return s.getLimits(ctx, date)
	}
	return []quik.Limit{}, s.err
}

func (s svc) GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	if s.getMoneyLimits != nil {
		return s.getMoneyLimits(ctx, date)
	}
	return []quik.MoneyLimit{}, s.err
}

func (s svc) GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	if s.getSecurityLimits != nil {
		return s.getSecurityLimits(ctx, date)
	}
	return []quik.SecurityLimit{}, s.err
}

func (s svc) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	if s.getSecurityLimitsOtc != nil {
		return s.getSecurityLimitsOtc(ctx, date)
	}
	return []quik.SecurityLimit{}, s.err
}

func (s svc) CreateMoneyLimit(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {
	if s.createMoneyLimit != nil {
		return s.createMoneyLimit(ctx, ml)
	}
	return ml, s.err
}

func (s svc) CreateSecurityLimit(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
	if s.createSecurityLimit != nil {
		return s.createSecurityLimit(ctx, sec)
	}
	return sec, s.err
}

func (s svc) CreateSecurityLimitOtc(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
	if s.createSecurityLimitOtc != nil {
		return s.createSecurityLimitOtc(ctx, sec)
	}
	return sec, s.err
}

func (s svc) GetPortfolio(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error) {
	if s.getPortfolio != nil {
		return s.getPortfolio(ctx, targetCcy)
	}
	return []quik.PortfolioEntry{}, s.err
}

const exampleURL = "http://example.ru/res"

func newTestHandler(s svc) *Handler {
	return NewHandler(handler.NewHandler(), s, zap.NewNop())
}

func reqWithQuery(t *testing.T, param string, value string) *http.Request {
	t.Helper()

	raw, err := url.Parse(exampleURL)
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}

	q := raw.Query()
	q.Set(param, value)
	raw.RawQuery = q.Encode()
	return httptest.NewRequest(http.MethodGet, raw.String(), nil)
}

func reqJSON(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, exampleURL, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	return req
}
