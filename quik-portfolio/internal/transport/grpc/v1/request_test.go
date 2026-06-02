package v1

import (
	"testing"
	"time"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	datepb "google.golang.org/genproto/googleapis/type/date"
)

func Test_extractReqFields(t *testing.T) {
	tests := []struct {
		name    string
		req     *quikv1.LimitsRequest
		want    limitsListQuery
		wantErr error
	}{
		{
			name: "все_поля_переданы",
			req: &quikv1.LimitsRequest{
				LoadDate:          &datepb.Date{Year: 2026, Month: 5, Day: 31},
				ClientCodes:       []string{"AA", "BB"},
				IncludeTotalCount: true,
				Limit:             25,
				Offset:            7,
			},
			want: limitsListQuery{
				Date:              time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local),
				Limit:             25,
				Offset:            7,
				ClientCodes:       []string{"AA", "BB"},
				IncludeTotalCount: true,
			},
		},
		{
			name: "нулевой_limit_заменяется_значением_по_умолчанию",
			req: &quikv1.LimitsRequest{
				LoadDate: &datepb.Date{Year: 2026, Month: 5, Day: 31},
			},
			want: limitsListQuery{
				Date:  time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local),
				Limit: defaultLimit,
			},
		},
		{
			name: "limit_больше_максимального_срезается_до_максимума",
			req: &quikv1.LimitsRequest{
				LoadDate: &datepb.Date{Year: 2026, Month: 5, Day: 31},
				Limit:    maxLimit + 1,
			},
			want: limitsListQuery{
				Date:  time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local),
				Limit: maxLimit,
			},
		},
		{
			name: "некорректная_дата",
			req: &quikv1.LimitsRequest{
				LoadDate: &datepb.Date{Year: 2026, Month: 2, Day: 31},
			},
			wantErr: md.ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractReqFields(tt.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
