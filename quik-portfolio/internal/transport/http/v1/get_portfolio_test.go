package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetPortfolio(t *testing.T) {
	t.Parallel()

	loadDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)
	quoteDate := time.Date(2025, 1, 3, 0, 0, 0, 0, time.Local)
	shortName := "Сбербанк"
	isin := "RU000A0JX0J2"

	tests := []struct {
		name       string
		req        *http.Request
		svc        svc
		wantBody   any
		wantDetail string
		wantErr    error
	}{
		{
			name: "валюта_по_умолчанию_rub",
			req:  httptest.NewRequest(http.MethodGet, exampleURL, nil),
			svc: svc{
				getPortfolio: func(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error) {
					return []quik.PortfolioEntry{
						{
							LimitType:      quik.LimitTypeSecurities,
							LoadDate:       loadDate,
							SourceDate:     sourceDate,
							ClientCode:     "AB12CD",
							FirmCode:       "COFE",
							FirmName:       "Фирма брокера",
							Instrument:     "SBER",
							ISIN:           &isin,
							ShortName:      &shortName,
							QuoteDate:      &quoteDate,
							Balance:        decimal.NewFromInt(10),
							MvCurrency:     "RUB",
							MvInCcy:        decimal.NewFromInt(1000),
							MvPrice:        decimal.NewFromInt(100),
							MvAccrued:      decimal.NewFromInt(0),
							MvTotal:        decimal.NewFromInt(1000),
							TargetCurrency: targetCcy,
						},
					}, nil
				},
			},
			wantBody: []portfolioEntryDTO{
				{
					LimitType:      "securities",
					LoadDate:       "2025-01-01",
					SourceDate:     "2025-01-02",
					ClientCode:     "AB12CD",
					FirmCode:       "COFE",
					FirmName:       "Фирма брокера",
					Instrument:     "SBER",
					ISIN:           "RU000A0JX0J2",
					ShortName:      "Сбербанк",
					QuoteDate:      "2025-01-03",
					QTY:            decimal.NewFromInt(10),
					MvCurrency:     "RUB",
					MvInCcy:        decimal.NewFromInt(1000),
					MvPrice:        decimal.NewFromInt(100),
					MvAccrued:      decimal.NewFromInt(0),
					MvTotal:        decimal.NewFromInt(1000),
					TargetCurrency: "RUB",
				},
			},
		},
		{
			name: "валюта_нормализуется",
			req:  httptest.NewRequest(http.MethodGet, exampleURL+"?targetCcy=+usd+", nil),
			svc: svc{
				getPortfolio: func(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error) {
					return []quik.PortfolioEntry{
						{
							LoadDate:       loadDate,
							SourceDate:     sourceDate,
							Balance:        decimal.NewFromInt(1),
							TargetCurrency: targetCcy,
						},
					}, nil
				},
			},
			wantBody: []portfolioEntryDTO{
				{
					LoadDate:       "2025-01-01",
					SourceDate:     "2025-01-02",
					QTY:            decimal.NewFromInt(1),
					TargetCurrency: "USD",
				},
			},
		},
		{
			name: "бизнес_ошибка",
			req:  httptest.NewRequest(http.MethodGet, exampleURL, nil),
			svc: svc{
				err: models.ErrBusinessValidation,
			},
			wantBody:   nil,
			wantDetail: models.ErrBusinessValidation.Error(),
			wantErr:    models.ErrBusinessValidation,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(tt.svc)
			body, detail, err := h.GetPortfolio(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
