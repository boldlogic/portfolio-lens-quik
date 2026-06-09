package v1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func jsonObjectKeys(t *testing.T, v any) map[string]json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	var out map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(b, &out))
	return out
}

func Test_moneyLimitDTO_jsonFieldNames(t *testing.T) {
	t.Parallel()
	dto := moneyLimitToDTO(quik.MoneyLimit{
		LoadDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		SourceDate:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
		ClientCode:   "AB12CD",
		CurrencyCode: "RUB",
		PositionCode: "EQTV",
		SettleCode:   quik.SettleCodeT2,
		FirmCode:     "NC0058900000",
		FirmName:     "Фирма брокера",
		Balance:      decimal.NewFromFloat(10.5),
	})
	keys := jsonObjectKeys(t, dto)
	for _, key := range []string{
		"loadDate", "sourceDate", "clientCode", "currency",
		"positionCode", "settleCode", "firmCode", "firmName", "balance",
	} {
		assert.Contains(t, keys, key)
	}
	for _, key := range []string{"currencyCode", "client_code"} {
		assert.NotContains(t, keys, key)
	}
}

func Test_securityLimitDTO_jsonFieldNames(t *testing.T) {
	t.Parallel()
	dto := securityLimitToDTO(quik.SecurityLimit{
		LoadDate:                time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		SourceDate:              time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
		ClientCode:              "AB12CD",
		SecCode:                 "SBER",
		TradeAccount:            "L01-00000F00",
		SettleCode:              quik.SettleCodeT2,
		FirmCode:                "NC0058900000",
		FirmName:                "Фирма брокера",
		Balance:                 decimal.NewFromFloat(10.5),
		AcquisitionCurrencyCode: "RUB",
		ISIN:                    "RU000A0JX0J2",
	})
	keys := jsonObjectKeys(t, dto)
	for _, key := range []string{
		"loadDate", "sourceDate", "clientCode", "secCode", "tradeAccount",
		"settleCode", "firmCode", "firmName", "balance", "acquisitionCurrencyCode",
	} {
		assert.Contains(t, keys, key)
	}
	for _, key := range []string{"ticker", "acquisition_ccy", "sec_code"} {
		assert.NotContains(t, keys, key)
	}
}
