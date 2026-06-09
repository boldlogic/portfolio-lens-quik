package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_moneyLimitDTO_jsonFieldNames(t *testing.T) {
	t.Parallel()
	keys := jsonObjectKeys(t, moneyLimitToDTO(fixtureMoneyLimit()))
	for _, key := range []string{
		"loadDate", "sourceDate", "clientCode", "currency",
		"positionCode", "settleCode", "firmCode", "firmName", "balance",
	} {
		assert.Contains(t, keys, key)
	}
	for _, key := range []string{"currencyCode", "client_code", "firm_code"} {
		assert.NotContains(t, keys, key)
	}
}

func Test_securityLimitDTO_jsonFieldNames(t *testing.T) {
	t.Parallel()
	keys := jsonObjectKeys(t, securityLimitToDTO(fixtureSecurityLimit()))
	for _, key := range []string{
		"loadDate", "sourceDate", "clientCode", "secCode", "tradeAccount",
		"settleCode", "firmCode", "firmName", "balance", "acquisitionCurrencyCode",
		"isin", "shortName",
	} {
		assert.Contains(t, keys, key)
	}
	for _, key := range []string{"ticker", "acquisition_ccy", "sec_code"} {
		assert.NotContains(t, keys, key)
	}
}

func Test_securityLimitDTO_jsonOmitempty(t *testing.T) {
	t.Parallel()
	in := fixtureSecurityLimit()
	in.ISIN = ""
	in.ShortName = ""
	keys := jsonObjectKeys(t, securityLimitToDTO(in))
	assert.NotContains(t, keys, "isin")
	assert.NotContains(t, keys, "shortName")
}

func Test_moneyLimitsResponseDTO_jsonFieldNames(t *testing.T) {
	t.Parallel()
	dto := moneyLimitsResponseDTO{
		Limits: []moneyLimitDTO{},
		Limit:  10,
		Offset: 5,
	}
	keys := jsonObjectKeys(t, dto)
	for _, key := range []string{"limits", "limit", "offset"} {
		assert.Contains(t, keys, key)
	}
}

func Test_securityLimitsResponseDTO_jsonFieldNames(t *testing.T) {
	t.Parallel()
	total := uint64(0)
	dto := securityLimitsResponseDTO{
		Limits:     []securityLimitDTO{},
		TotalCount: &total,
		Limit:      10,
		Offset:     5,
	}
	keys := jsonObjectKeys(t, dto)
	for _, key := range []string{"limits", "totalCount", "limit", "offset"} {
		assert.Contains(t, keys, key)
	}
}
