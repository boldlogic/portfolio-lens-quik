package v1

import (
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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

func Test_securityLimitToDTO(t *testing.T) {
	t.Parallel()

	isin := "RU000A0JX0J2"

	tests := []struct {
		name string
		in   quik.SecurityLimit
		want securityLimitDTO
	}{
		{
			name: "все_поля_с_isin",
			in: quik.SecurityLimit{
				LoadDate:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
				SourceDate:     time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
				ClientCode:     "AB12CD",
				Ticker:         "SBER",
				TradeAccount:   "L01-00000F00",
				SettleCode:     quik.SettleCodeT2,
				FirmCode:       "COFE",
				FirmName:       "Фирма брокера",
				Balance:        decimal.RequireFromString("15.25"),
				AcquisitionCcy: "RUB",
				ISIN:           &isin,
			},
			want: securityLimitDTO{
				LoadDate:       "2025-01-01",
				SourceDate:     "2025-01-02",
				ClientCode:     "AB12CD",
				Ticker:         "SBER",
				TradeAccount:   "L01-00000F00",
				SettleCode:     "T2",
				FirmCode:       "COFE",
				FirmName:       "Фирма брокера",
				Balance:        decimal.RequireFromString("15.25"),
				AcquisitionCcy: "RUB",
				ISIN:           "RU000A0JX0J2",
			},
		},
		{
			name: "nil_isin_оставляет_пустую_строку",
			in: quik.SecurityLimit{
				LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
				Balance:    decimal.RequireFromString("-10.50"),
			},
			want: securityLimitDTO{
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
			got := securityLimitToDTO(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_securityLimitsToResp(t *testing.T) {
	t.Parallel()

	isin := "RU000A0JX0J2"

	tests := []struct {
		name string
		in   []quik.SecurityLimit
		want []securityLimitDTO
	}{
		{
			name: "есть_лимиты",
			in: []quik.SecurityLimit{
				{
					LoadDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
					SourceDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local),
					Ticker:     "SBER",
					Balance:    decimal.NewFromInt(10),
					ISIN:       &isin,
				},
				{
					LoadDate:   time.Date(2026, 12, 30, 0, 0, 0, 0, time.Local),
					SourceDate: time.Date(2026, 12, 30, 0, 0, 0, 0, time.Local),
					Ticker:     "GAZP",
					Balance:    decimal.RequireFromString("-1.50"),
				},
			},
			want: []securityLimitDTO{
				{
					LoadDate:   "2026-12-31",
					SourceDate: "2026-12-31",
					Ticker:     "SBER",
					Balance:    decimal.NewFromInt(10),
					ISIN:       "RU000A0JX0J2",
				},
				{
					LoadDate:   "2026-12-30",
					SourceDate: "2026-12-30",
					Ticker:     "GAZP",
					Balance:    decimal.RequireFromString("-1.50"),
				},
			},
		},
		{
			name: "nil_слайс_возвращает_пустой_ответ",
			in:   nil,
			want: []securityLimitDTO{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := securityLimitsToResp(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
