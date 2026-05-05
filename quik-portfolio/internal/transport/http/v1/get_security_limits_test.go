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

func TestGetSecurityLimits(t *testing.T) {
	t.Parallel()

	isin := "RU000A0JX0J2"
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
			req:  reqWithQuery(t, "date", "2025-01-01"),
			svc: svc{
				getSecurityLimits: func(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
					return []quik.SecurityLimit{
						{
							LoadDate:       date,
							SourceDate:     sourceDate,
							ClientCode:     "AB12CD",
							Ticker:         "SBER",
							TradeAccount:   "L01-00000F00",
							SettleCode:     quik.SettleCodeT2,
							FirmCode:       "COFE",
							FirmName:       "Фирма брокера",
							Balance:        decimal.RequireFromString("10.5"),
							AcquisitionCcy: "RUB",
							ISIN:           &isin,
						},
					}, nil
				},
			},
			wantBody: []securityLimitDTO{
				{
					LoadDate:       "2025-01-01",
					SourceDate:     "2025-01-02",
					ClientCode:     "AB12CD",
					Ticker:         "SBER",
					TradeAccount:   "L01-00000F00",
					SettleCode:     "T2",
					FirmCode:       "COFE",
					FirmName:       "Фирма брокера",
					Balance:        decimal.RequireFromString("10.5"),
					AcquisitionCcy: "RUB",
					ISIN:           "RU000A0JX0J2",
				},
			},
		},
		{
			name:       "некорректная_дата",
			req:        reqWithQuery(t, "date", "2025-01"),
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
			body, detail, err := h.GetSecurityLimits(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
