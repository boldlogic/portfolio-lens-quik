package v1

import (
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
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
				fixtureMoneyLimit(),
			},
			limit:             25,
			offset:            7,
			totalCount:        &totalCount,
			includeTotalCount: true,
			want: moneyLimitsResponseDTO{
				Limits: []moneyLimitDTO{
					fixtureMoneyLimitDTO(),
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
