package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateMoneyLimit(t *testing.T) {
	t.Parallel()

	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name       string
		req        *http.Request
		svc        svc
		wantBody   any
		wantDetail string
		wantErr    error
	}{
		{
			name: "успешный_запрос",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"currency":"RUB",
				"positionCode":"EQTV",
				"settleCode":"Tx",
				"firmName":"Фирма брокера",
				"balance":331.10
			}`),
			svc: svc{
				createMoneyLimit: func(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {
					ml.SourceDate = sourceDate
					ml.FirmCode = "COFE"
					return ml, nil
				},
			},
			wantBody: moneyLimitDTO{
				LoadDate:     "2025-01-01",
				SourceDate:   "2025-01-02",
				ClientCode:   "AB12CD",
				Currency:     "RUB",
				PositionCode: "EQTV",
				SettleCode:   "Tx",
				FirmCode:     "COFE",
				FirmName:     "Фирма брокера",
				Balance:      decimal.NewFromFloat(331.10),
			},
		},
		{
			name:       "некорректная_дата",
			req:        reqJSON(`{"loadDate":"2025-01","clientCode":"AB12CD","currency":"RUB","firmName":"Фирма брокера"}`),
			wantBody:   nil,
			wantDetail: dates.ErrWrongDateFormat.Error(),
			wantErr:    models.ErrValidation,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(tt.svc)
			body, detail, err := h.CreateMoneyLimit(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCreateSecurityLimit(t *testing.T) {
	t.Parallel()

	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name       string
		req        *http.Request
		svc        svc
		wantBody   any
		wantDetail string
		wantErr    error
	}{
		{
			name: "успешный_запрос",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"settleCode":"T2",
				"firmName":"Фирма брокера",
				"balance":10.5,
				"acquisitionCcy":"RUB",
				"isin":"RU000A0JX0J2"
			}`),
			svc: svc{
				createSecurityLimit: func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
					sec.SourceDate = sourceDate
					sec.FirmCode = "COFE"
					return sec, nil
				},
			},
			wantBody: securityLimitDTO{
				LoadDate:       "2025-01-01",
				SourceDate:     "2025-01-02",
				ClientCode:     "AB12CD",
				Ticker:         "SBER",
				TradeAccount:   "L01-00000F00",
				SettleCode:     "T2",
				FirmCode:       "COFE",
				FirmName:       "Фирма брокера",
				Balance:        decimal.NewFromFloat(10.5),
				AcquisitionCcy: "RUB",
				ISIN:           "RU000A0JX0J2",
			},
		},
		{
			name: "конфликт_ключа",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"firmName":"Фирма брокера",
				"balance":10.5
			}`),
			svc: svc{
				err: models.ErrConflict,
			},
			wantBody:   nil,
			wantDetail: models.ErrConflict.Error(),
			wantErr:    models.ErrConflict,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(tt.svc)
			body, detail, err := h.CreateSecurityLimit(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCreateSecurityLimitOtc(t *testing.T) {
	t.Parallel()

	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name       string
		req        *http.Request
		svc        svc
		wantBody   any
		wantDetail string
		wantErr    error
	}{
		{
			name: "успешный_запрос_без_trade_account",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"OTC_BOND",
				"settleCode":"T0",
				"firmName":"Фирма брокера",
				"balance":2,
				"acquisitionCcy":"USD"
			}`),
			svc: svc{
				createSecurityLimitOtc: func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
					sec.SourceDate = sourceDate
					sec.FirmCode = "COFE"
					return sec, nil
				},
			},
			wantBody: securityLimitDTO{
				LoadDate:       "2025-01-01",
				SourceDate:     "2025-01-02",
				ClientCode:     "AB12CD",
				Ticker:         "OTC_BOND",
				SettleCode:     "T0",
				FirmCode:       "COFE",
				FirmName:       "Фирма брокера",
				Balance:        decimal.NewFromFloat(2),
				AcquisitionCcy: "USD",
			},
		},
		{
			name:       "битый_json",
			req:        reqJSON(`{"loadDate":"2025-01-01"`),
			wantBody:   nil,
			wantDetail: "unexpected EOF",
			wantErr:    models.ErrValidation,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(tt.svc)
			body, detail, err := h.CreateSecurityLimitOtc(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
