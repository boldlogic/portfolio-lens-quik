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
				"firmCode": "NC0058900000",
				"balance":"331.10"
			}`),
			svc: svc{
				createMoneyLimit: func(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {
					ml.SourceDate = sourceDate
					ml.FirmCode = "NC0058900000"
					ml.FirmName = "Фирма брокера"
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
				FirmCode:     "NC0058900000",
				FirmName:     "Фирма брокера",
				Balance:      decimal.RequireFromString("331.10"),
			},
		},
		{
			name:       "некорректная_дата",
			req:        reqJSON(`{"loadDate":"2025-01","clientCode":"AB12CD","currency":"RUB","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: dates.ErrWrongDateFormat.Error(),
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_clientCode",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"","currency":"RUB","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "clientCode>12",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"1234567890123","currency":"RUB","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_currency",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"TBANK","currency":"","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: "currency",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_firmCode",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"TBANK","currency":"RUR","firmCode": ""}`),
			wantBody:   nil,
			wantDetail: "firmCode",
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
