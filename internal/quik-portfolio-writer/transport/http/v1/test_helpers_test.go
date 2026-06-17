package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type svc struct {
	err                    error
	createMoneyLimit       func(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error)
	createSecurityLimit    func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error)
	createSecurityLimitOtc func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error)
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

const exampleURL = "http://example.ru/res"

func newTestHandler(s svc) *Handler {
	return NewHandler(handler.NewHandler(), s, zap.NewNop())
}

func reqJSON(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, exampleURL, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	return req
}
