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

	firmNameEmpty := fixtureMoneyLimit()
	firmNameEmpty.FirmName = ""
	wantFirmNameEmpty := fixtureMoneyLimitDTO()
	wantFirmNameEmpty.FirmName = ""

	currencyRUR := fixtureMoneyLimit()
	currencyRUR.CurrencyCode = "RUR"
	wantCurrencyRUB := fixtureMoneyLimitDTO()
	wantCurrencyRUB.Currency = "RUB"

	tests := []struct {
		name string
		in   quik.MoneyLimit
		want moneyLimitDTO
	}{
		{
			name: "все_поля_заполнены",
			in:   fixtureMoneyLimit(),
			want: fixtureMoneyLimitDTO(),
		},
		{
			name: "firmName_пустой",
			in:   firmNameEmpty,
			want: wantFirmNameEmpty,
		},
		{
			name: "currency_rur_из_БД_rub_в_dto",
			in:   currencyRUR,
			want: wantCurrencyRUB,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, moneyLimitToDTO(tt.in))
		})
	}
}

func Test_securityLimitToDTO(t *testing.T) {
	t.Parallel()

	isinEmpty := fixtureSecurityLimit()
	isinEmpty.ISIN = ""
	wantIsinEmpty := fixtureSecurityLimitDTO()
	wantIsinEmpty.ISIN = ""

	firmNameEmpty := fixtureSecurityLimit()
	firmNameEmpty.FirmName = ""
	wantFirmNameEmpty := fixtureSecurityLimitDTO()
	wantFirmNameEmpty.FirmName = ""

	shortNameTrim := fixtureSecurityLimit()
	shortNameTrim.ShortName = " Сбербанк "
	wantShortNameTrim := fixtureSecurityLimitDTO()
	wantShortNameTrim.ShortName = "Сбербанк"

	acqRUR := fixtureSecurityLimit()
	acqRUR.AcquisitionCurrencyCode = "RUR"
	wantAcqRUB := fixtureSecurityLimitDTO()
	wantAcqRUB.AcquisitionCurrencyCode = "RUB"

	tests := []struct {
		name string
		in   quik.SecurityLimit
		want securityLimitDTO
	}{
		{
			name: "все_поля_заполнены",
			in:   fixtureSecurityLimit(),
			want: fixtureSecurityLimitDTO(),
		},
		{
			name: "isin_пустой",
			in:   isinEmpty,
			want: wantIsinEmpty,
		},
		{
			name: "firmName_пустой",
			in:   firmNameEmpty,
			want: wantFirmNameEmpty,
		},
		{
			name: "shortName_trim",
			in:   shortNameTrim,
			want: wantShortNameTrim,
		},
		{
			name: "acquisitionCurrency_rur_из_БД_rub_в_dto",
			in:   acqRUR,
			want: wantAcqRUB,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, securityLimitToDTO(tt.in))
		})
	}
}

func Test_securityLimitsToDTO(t *testing.T) {
	t.Parallel()

	second := fixtureSecurityLimit()
	second.LoadDate = time.Date(2026, 5, 30, 0, 0, 0, 0, time.Local)
	second.SourceDate = time.Date(2026, 5, 29, 0, 0, 0, 0, time.Local)
	second.SecCode = fixtureSecCodeGAZP
	second.Balance = decimal.RequireFromString("5.50")
	second.ISIN = ""
	second.ShortName = ""

	wantSecond := fixtureSecurityLimitDTO()
	wantSecond.LoadDate = "2026-05-30"
	wantSecond.SourceDate = "2026-05-29"
	wantSecond.SecCode = fixtureSecCodeGAZP
	wantSecond.Balance = decimal.RequireFromString("5.50")
	wantSecond.ISIN = ""
	wantSecond.ShortName = ""

	tests := []struct {
		name string
		in   []quik.SecurityLimit
		want []securityLimitDTO
	}{
		{
			name: "один_элемент",
			in:   []quik.SecurityLimit{fixtureSecurityLimit()},
			want: []securityLimitDTO{fixtureSecurityLimitDTO()},
		},
		{
			name: "два_элемента_сохраняют_порядок",
			in:   []quik.SecurityLimit{fixtureSecurityLimit(), second},
			want: []securityLimitDTO{fixtureSecurityLimitDTO(), wantSecond},
		},
		{
			name: "nil_слайс_возвращает_пустой_слайс",
			in:   nil,
			want: []securityLimitDTO{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, securityLimitsToDTO(tt.in))
		})
	}
}
