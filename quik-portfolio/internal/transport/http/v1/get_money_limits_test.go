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
			name: "успешный_запрос",
			svc: svc{
				getMoneyLimits: func(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
					return []quik.MoneyLimit{
						{
							LoadDate:     date,
							SourceDate:   sourceDate,
							ClientCode:   "AB12CD",
							Currency:     "RUB",
							PositionCode: "EQTV",
							SettleCode:   quik.SettleCodeTx,
							FirmCode:     "COFE",
							FirmName:     "Фирма брокера",
							Balance:      decimal.RequireFromString("331.10"),
						},
					}, nil
				},
			},
			wantBody: []moneyLimitDTO{
				{
					LoadDate:     "2025-01-01",
					SourceDate:   "2025-01-02",
					ClientCode:   "AB12CD",
					Currency:     "RUB",
					PositionCode: "EQTV",
					SettleCode:   "Tx",
					FirmCode:     "COFE",
					FirmName:     "Фирма брокера",
					Balance:      decimal.RequireFromString("331.10"),
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
			body, detail, err := h.GetMoneyLimits(reqWithQuery(t, "date", "2025-01-01"))
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_moneyLimitToDTO(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   quik.MoneyLimit
		want moneyLimitDTO
	}{
		{
			name: "все_поля",
			in: quik.MoneyLimit{
				LoadDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
				SourceDate:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
				ClientCode:   "AB12CD",
				Currency:     "RUB",
				PositionCode: "EQTV",
				SettleCode:   quik.SettleCodeTx,
				FirmName:     "Фирма брокера",
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
				FirmName:     "Фирма брокера",
				Balance:      decimal.NewFromInt(331),
			},
		},
		{
			name: "отрицательный_баланс",
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
			name: "дробный_баланс",
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
			name: "отрицательный_дробный_баланс",
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
			name: "есть_лимиты",
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
			name: "nil_слайс_возвращает_пустой_ответ",
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
