package v1

import (
	"testing"
	"time"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	datepb "google.golang.org/genproto/googleapis/type/date"
)

func Test_protoDateToTime(t *testing.T) {
	tests := []struct {
		name    string
		d       *datepb.Date
		want    time.Time
		wantErr error
	}{
		{
			name: "минимальная_допустимая_дата",
			d:    &datepb.Date{Year: 2000, Month: 1, Day: 1},
			want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			name: "максимальная_допустимая_дата",
			d:    &datepb.Date{Year: 2999, Month: 12, Day: 31},
			want: time.Date(2999, 12, 31, 0, 0, 0, 0, time.Local),
		},
		{
			name:    "год_меньше_минимального",
			d:       &datepb.Date{Year: 1999, Month: 1, Day: 1},
			wantErr: md.ErrValidation,
		},
		{
			name:    "год_больше_максимального",
			d:       &datepb.Date{Year: 3000, Month: 1, Day: 1},
			wantErr: md.ErrValidation,
		},
		{
			name:    "месяц_меньше_минимального",
			d:       &datepb.Date{Year: 2026, Month: 0, Day: 1},
			wantErr: md.ErrValidation,
		},
		{
			name:    "месяц_больше_максимального",
			d:       &datepb.Date{Year: 2026, Month: 13, Day: 1},
			wantErr: md.ErrValidation,
		},
		{
			name:    "день_меньше_минимального",
			d:       &datepb.Date{Year: 2026, Month: 1, Day: 0},
			wantErr: md.ErrValidation,
		},
		{
			name:    "день_больше_максимального",
			d:       &datepb.Date{Year: 2026, Month: 1, Day: 32},
			wantErr: md.ErrValidation,
		},
		{
			name:    "день_не_существует_в_месяце",
			d:       &datepb.Date{Year: 2026, Month: 2, Day: 31},
			wantErr: md.ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := protoDateToTime(tt.d)
			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
