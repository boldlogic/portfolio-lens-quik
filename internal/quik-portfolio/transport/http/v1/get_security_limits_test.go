package v1

import (
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_securityLimitsToResponseDTO(t *testing.T) {
	totalCount := uint64(1)
	zeroTotalCount := uint64(0)

	shortNameInput := fixtureSecurityLimit()
	shortNameInput.ShortName = " Сбербанк "
	shortNameWant := fixtureSecurityLimitDTO()
	shortNameWant.ShortName = "Сбербанк"

	tests := []struct {
		name              string
		sls               []quik.SecurityLimit
		limit             uint32
		offset            uint64
		totalCount        *uint64
		includeTotalCount bool
		want              securityLimitsResponseDTO
	}{
		{
			name:   "маппит_один_лимит",
			sls:    []quik.SecurityLimit{fixtureSecurityLimit()},
			limit:  25,
			offset: 7,
			want: securityLimitsResponseDTO{
				Limits: []securityLimitDTO{fixtureSecurityLimitDTO()},
				Limit:  25,
				Offset: 7,
			},
		},
		{
			name:   "shortName_обрезает_пробелы",
			sls:    []quik.SecurityLimit{shortNameInput},
			limit:  10,
			want: securityLimitsResponseDTO{
				Limits: []securityLimitDTO{shortNameWant},
				Limit:  10,
			},
		},
		{
			name:              "прокидывает_totalCount",
			limit:             10,
			totalCount:        &totalCount,
			includeTotalCount: true,
			want: securityLimitsResponseDTO{
				Limits:     []securityLimitDTO{},
				TotalCount: &totalCount,
				Limit:      10,
			},
		},
		{
			name:              "нет_лимитов_и_totalCount_0",
			limit:             10,
			totalCount:        &zeroTotalCount,
			includeTotalCount: true,
			want: securityLimitsResponseDTO{
				Limits:     []securityLimitDTO{},
				TotalCount: &zeroTotalCount,
				Limit:      10,
			},
		},
		{
			name:   "totalCount_не_запрошен",
			limit:  10,
			offset: 5,
			want: securityLimitsResponseDTO{
				Limits: []securityLimitDTO{},
				Limit:  10,
				Offset: 5,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := securityLimitsToResponseDTO(tt.sls, tt.limit, tt.offset, tt.totalCount, tt.includeTotalCount)
			if tt.includeTotalCount {
				require.NotNil(t, got.TotalCount)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
