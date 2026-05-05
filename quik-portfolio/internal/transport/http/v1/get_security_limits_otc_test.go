package v1

import (
	"context"
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetSecurityLimitsOtc(t *testing.T) {
	t.Parallel()

	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name       string
		svc        svc
		wantBody   any
		wantDetail string
		wantErr    error
	}{
		{
			name: "успешный_запрос",
			svc: svc{
				getSecurityLimitsOtc: func(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
					return []quik.SecurityLimit{
						{
							LoadDate:       date,
							SourceDate:     sourceDate,
							ClientCode:     "AB12CD",
							Ticker:         "OTC_BOND",
							SettleCode:     quik.SettleCodeT0,
							FirmCode:       "COFE",
							FirmName:       "Фирма брокера",
							Balance:        decimal.NewFromInt(2),
							AcquisitionCcy: "USD",
						},
					}, nil
				},
			},
			wantBody: []securityLimitDTO{
				{
					LoadDate:       "2025-01-01",
					SourceDate:     "2025-01-02",
					ClientCode:     "AB12CD",
					Ticker:         "OTC_BOND",
					SettleCode:     "T0",
					FirmCode:       "COFE",
					FirmName:       "Фирма брокера",
					Balance:        decimal.NewFromInt(2),
					AcquisitionCcy: "USD",
				},
			},
		},
		{
			name: "бизнес_ошибка",
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
			body, detail, err := h.GetSecurityLimitsOtc(reqWithQuery(t, "date", "2025-01-01"))
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
