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
	"github.com/stretchr/testify/require"
)

func Test_normalizeClientCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want []string
		err  bool
	}{
		{name: "empty", raw: "", want: nil},
		{name: "commas_only", raw: ",,,", want: nil},
		{name: "trim_upper_skip_empty", raw: " ab1 , , Ab2 ", want: []string{"AB1", "AB2"}},
		{name: "single_code", raw: "X1", want: []string{"X1"}},
		{name: "too_long", raw: "1234567890123", err: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := normalizeClientCodes(tt.raw)
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetMoneyLimits(t *testing.T) {
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
			name: "success",
			svc: svc{
				getMoneyLimitsWithFilters: func(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) ([]quik.MoneyLimit, int, error) {
					assert.Equal(t, 50, limit)
					assert.Equal(t, 0, offset)
					assert.Nil(t, clientCodes)
					return []quik.MoneyLimit{
						{
							LoadDate:     date,
							SourceDate:   sourceDate,
							ClientCode:   "AB12CD",
							Currency:     "RUB",
							PositionCode: "EQTV",
							SettleCode:   quik.SettleCodeTx,
							FirmCode:     "COFE",
							FirmName:     "Broker firm",
							Balance:      decimal.RequireFromString("331.10"),
						},
					}, 42, nil
				},
			},
			wantBody: moneyLimitsDTO{
				Limits: []moneyLimitDTO{
					{
						LoadDate:     "2025-01-01",
						SourceDate:   "2025-01-02",
						ClientCode:   "AB12CD",
						Currency:     "RUB",
						PositionCode: "EQTV",
						SettleCode:   "Tx",
						FirmCode:     "COFE",
						FirmName:     "Broker firm",
						Balance:      decimal.RequireFromString("331.10"),
					},
				},
				TotalCount: 42,
				Limit:      50,
				Offset:     0,
			},
		},
		{
			name: "business_error",
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
			body, detail, err := h.GetMoneyLimits(reqWithQuery(t, "date", "2025-01-01"))
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetMoneyLimits_passesPaginationAndClientCodes(t *testing.T) {
	t.Parallel()

	h := newTestHandler(svc{
		getMoneyLimitsWithFilters: func(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) ([]quik.MoneyLimit, int, error) {
			assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local), date)
			assert.Equal(t, 25, limit)
			assert.Equal(t, 75, offset)
			assert.Equal(t, []string{"AB1", "CD2"}, clientCodes)
			return nil, 123, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "http://example.ru/res?date=2025-01-01&limit=25&offset=75&clientCodes=ab1,%20CD2", nil)
	body, detail, err := h.GetMoneyLimits(req)

	require.NoError(t, err)
	assert.Empty(t, detail)
	assert.Equal(t, moneyLimitsDTO{
		Limits:     []moneyLimitDTO{},
		TotalCount: 123,
		Limit:      25,
		Offset:     75,
	}, body)
}

func Test_moneyLimitToDTO(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   quik.MoneyLimit
		want moneyLimitDTO
	}{
		{
			name: "all_fields",
			in: quik.MoneyLimit{
				LoadDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
				SourceDate:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
				ClientCode:   "AB12CD",
				Currency:     "RUB",
				PositionCode: "EQTV",
				SettleCode:   quik.SettleCodeTx,
				FirmName:     "Broker firm",
				FirmCode:     "COFE",
				Balance:      decimal.NewFromInt(331),
			},
			want: moneyLimitDTO{
				LoadDate:     "2025-01-01",
				SourceDate:   "2025-01-02",
				ClientCode:   "AB12CD",
				Currency:     "RUB",
				PositionCode: "EQTV",
				SettleCode:   "Tx",
				FirmCode:     "COFE",
				FirmName:     "Broker firm",
				Balance:      decimal.NewFromInt(331),
			},
		},
		{
			name: "negative_balance",
			in: quik.MoneyLimit{
				LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				Balance:    decimal.NewFromInt(-331),
			},
			want: moneyLimitDTO{
				LoadDate:   "2026-12-31",
				SourceDate: "2026-12-31",
				Balance:    decimal.NewFromInt(-331),
			},
		},
		{
			name: "fractional_balance",
			in: quik.MoneyLimit{
				LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				Balance:    decimal.RequireFromString("11.59"),
			},
			want: moneyLimitDTO{
				LoadDate:   "2026-12-31",
				SourceDate: "2026-12-31",
				Balance:    decimal.RequireFromString("11.59"),
			},
		},
		{
			name: "negative_fractional_balance",
			in: quik.MoneyLimit{
				LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				Balance:    decimal.RequireFromString("-10.50"),
			},
			want: moneyLimitDTO{
				LoadDate:   "2026-12-31",
				SourceDate: "2026-12-31",
				Balance:    decimal.RequireFromString("-10.50"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := moneyLimitToDTO(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_moneyLimitsToResp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []quik.MoneyLimit
		want []moneyLimitDTO
	}{
		{
			name: "has_limits",
			in: []quik.MoneyLimit{
				{
					LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
					SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
					Balance:    decimal.NewFromInt(-331),
				},
				{
					LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
					SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
					Balance:    decimal.RequireFromString("-10.50"),
				},
			},
			want: []moneyLimitDTO{
				{
					LoadDate:   "2026-12-31",
					SourceDate: "2026-12-31",
					Balance:    decimal.NewFromInt(-331),
				},
				{
					LoadDate:   "2026-12-31",
					SourceDate: "2026-12-31",
					Balance:    decimal.RequireFromString("-10.50"),
				},
			},
		},
		{
			name: "nil_slice_returns_empty_response",
			in:   nil,
			want: []moneyLimitDTO{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := moneyLimitsToResp(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
