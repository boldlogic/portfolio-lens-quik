package v1

import (
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_moneyLimitsToResponseDTO(t *testing.T) {
	totalCount := uint64(1)
	zeroTotalCount := uint64(0)

	tests := []struct {
		name              string
		mls               []quik.MoneyLimit
		limit             uint32
		offset            uint64
		totalCount        *uint64
		includeTotalCount bool
		want              moneyLimitsResponseDTO
	}{
		{
			name: "есть_лимиты_и_totalCount",
			mls: []quik.MoneyLimit{
				{
					LoadDate:     time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local),
					SourceDate:   time.Date(2026, 5, 30, 0, 0, 0, 0, time.Local),
					ClientCode:   "AA",
					Currency:     "RUB",
					PositionCode: "EQTV",
					SettleCode:   quik.SettleCodeT2,
					FirmCode:     "COFE",
					FirmName:     "Брокер",
					Balance:      decimal.RequireFromString("10.25"),
				},
			},
			limit:             25,
			offset:            7,
			totalCount:        &totalCount,
			includeTotalCount: true,
			want: moneyLimitsResponseDTO{
				Limits: []moneyLimitDTO{
					{
						LoadDate:     "2026-05-31",
						SourceDate:   "2026-05-30",
						ClientCode:   "AA",
						Currency:     "RUB",
						PositionCode: "EQTV",
						SettleCode:   "T2",
						FirmCode:     "COFE",
						FirmName:     "Брокер",
						Balance:      decimal.RequireFromString("10.25"),
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
			want: moneyLimitsResponseDTO{
				Limits:     []moneyLimitDTO{},
				TotalCount: &zeroTotalCount,
				Limit:      10,
			},
		},
		{
			name:   "totalCount_не_запрошен",
			limit:  10,
			offset: 5,
			want: moneyLimitsResponseDTO{
				Limits: []moneyLimitDTO{},
				Limit:  10,
				Offset: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := moneyLimitsToResponseDTO(tt.mls, tt.limit, tt.offset, tt.totalCount, tt.includeTotalCount)
			if tt.includeTotalCount {
				require.NotNil(t, got.TotalCount)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
