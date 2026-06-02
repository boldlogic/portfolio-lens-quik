package v1

import (
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_securityLimitsWithPaginationToResp(t *testing.T) {
	totalCount := uint64(2)
	zeroTotalCount := uint64(0)
	isin := "RU000A0JX0J2"

	tests := []struct {
		name              string
		sls               []quik.SecurityLimit
		limit             uint32
		offset            uint64
		totalCount        *uint64
		includeTotalCount bool
		want              securityLimitsDTO
	}{
		{
			name: "есть_лимиты_и_totalCount",
			sls: []quik.SecurityLimit{
				{
					LoadDate:       time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local),
					SourceDate:     time.Date(2026, 5, 30, 0, 0, 0, 0, time.Local),
					ClientCode:     "AA",
					Ticker:         "SBER",
					TradeAccount:   "L01-00000F00",
					SettleCode:     quik.SettleCodeT2,
					FirmCode:       "COFE",
					FirmName:       "Брокер",
					Balance:        decimal.RequireFromString("10.25"),
					AcquisitionCcy: "RUB",
					ISIN:           &isin,
					ShortName:      " Сбербанк ",
				},
				{
					LoadDate:       time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local),
					SourceDate:     time.Date(2026, 5, 30, 0, 0, 0, 0, time.Local),
					ClientCode:     "AA",
					Ticker:         "GAZP",
					TradeAccount:   "L01-00000F00",
					SettleCode:     quik.SettleCodeTx,
					FirmCode:       "COFE",
					FirmName:       "Брокер",
					Balance:        decimal.RequireFromString("5.50"),
					AcquisitionCcy: "RUB",
				},
			},
			limit:             25,
			offset:            7,
			totalCount:        &totalCount,
			includeTotalCount: true,
			want: securityLimitsDTO{
				Limits: []securityLimitDTO{
					{
						LoadDate:       "2026-05-31",
						SourceDate:     "2026-05-30",
						ClientCode:     "AA",
						Ticker:         "SBER",
						TradeAccount:   "L01-00000F00",
						SettleCode:     "T2",
						FirmCode:       "COFE",
						FirmName:       "Брокер",
						Balance:        decimal.RequireFromString("10.25"),
						AcquisitionCcy: "RUB",
						ISIN:           "RU000A0JX0J2",
						ShortName:      "Сбербанк",
					},
					{
						LoadDate:       "2026-05-31",
						SourceDate:     "2026-05-30",
						ClientCode:     "AA",
						Ticker:         "GAZP",
						TradeAccount:   "L01-00000F00",
						SettleCode:     "Tx",
						FirmCode:       "COFE",
						FirmName:       "Брокер",
						Balance:        decimal.RequireFromString("5.50"),
						AcquisitionCcy: "RUB",
					},
				},
				TotalCount: &totalCount,
				Limit:      25,
				Offset:     7,
			},
		},
		{
			name:              "нет_лимитов_и_totalCount_0",
			limit:             10,
			offset:            0,
			totalCount:        &zeroTotalCount,
			includeTotalCount: true,
			want: securityLimitsDTO{
				Limits:     []securityLimitDTO{},
				TotalCount: &zeroTotalCount,
				Limit:      10,
			},
		},
		{
			name:   "totalCount_не_запрошен",
			limit:  10,
			offset: 5,
			want: securityLimitsDTO{
				Limits: []securityLimitDTO{},
				Limit:  10,
				Offset: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := securityLimitsWithPaginationToResp(tt.sls, tt.limit, tt.offset, tt.totalCount, tt.includeTotalCount)
			if tt.includeTotalCount {
				require.NotNil(t, got.TotalCount)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
