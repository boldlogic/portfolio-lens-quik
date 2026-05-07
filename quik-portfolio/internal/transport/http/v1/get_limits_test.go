package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_extractDateQueryParam(t *testing.T) {
	t.Parallel()
	h := newTestHandler(svc{})

	tests := []struct {
		name    string
		inParam string
		inValue string
		want    time.Time
		wantErr error
	}{
		{
			name:    "корректная_дата",
			inParam: "date",
			inValue: "2025-01-01",
			want:    time.Date(2025, 01, 01, 0, 0, 0, 0, time.Local),
			wantErr: nil,
		},
		{
			name:    "дата_по_умолчанию",
			inParam: "",
			inValue: "",
			want:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
			wantErr: nil,
		},
		{
			name:    "некорректная_дата",
			inParam: "date",
			inValue: "2025-01",
			want:    time.Time{},
			wantErr: dates.ErrWrongDateFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := url.Parse(exampleURL)
			if err != nil {
				t.Fatalf("ParseListPagination: %v", err)
			}
			q := raw.Query()
			q.Set(tt.inParam, tt.inValue)
			raw.RawQuery = q.Encode()
			req := httptest.NewRequest("GET", raw.String(), nil)

			got, err := h.extractDateQueryParam(req)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetLimits(t *testing.T) {
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
				getLimits: func(ctx context.Context, date time.Time) ([]quik.Limit, error) {
					return []quik.Limit{
						{
							LimitType:      quik.LimitTypeSecurities,
							LoadDate:       date,
							SourceDate:     sourceDate,
							ClientCode:     "AB12CD",
							InstrumentCode: "SBER",
							ISIN:           &isin,
							SettleCode:     quik.SettleCodeT2,
							FirmCode:       "COFE",
							FirmName:       "Фирма брокера",
							Balance:        decimal.RequireFromString("10.5"),
							AcquisitionCcy: "RUB",
						},
					}, nil
				},
			},
			wantBody: []limitDTO{
				{
					LimitType:      "securities",
					LoadDate:       "2025-01-01",
					SourceDate:     "2025-01-02",
					ClientCode:     "AB12CD",
					Instrument:     "SBER",
					ISIN:           "RU000A0JX0J2",
					SettleCode:     "T2",
					FirmCode:       "COFE",
					FirmName:       "Фирма брокера",
					Balance:        decimal.RequireFromString("10.5"),
					AcquisitionCcy: "RUB",
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
			body, detail, err := h.GetLimits(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
