package quik

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseCurrencyCode(t *testing.T) {
	tests := []struct {
		name    string
		rawCode string
		want    CurrencyCode
		wantErr error
	}{
		{
			name:    "обычная_валюта",
			rawCode: "CNY",
			want:    "CNY",
		},
		{
			name:    "RUB",
			rawCode: "RUB",
			want:    "RUB",
		},
		{
			name:    "USD",
			rawCode: "USD",
			want:    "USD",
		},
		{
			name:    "GLD-XAU",
			rawCode: "GLD",
			want:    "XAU",
		},
		{
			name:    "SUR-RUB",
			rawCode: "SUR",
			want:    "RUB",
		},
		{
			name:    "RUR-RUB",
			rawCode: "RUR",
			want:    "RUB",
		},
		{
			name:    "USDX-USD",
			rawCode: "USDX",
			want:    "USD",
		},
		{
			name:    "SLV-XAG",
			rawCode: "SLV",
			want:    "XAG",
		},
		{
			name:    "PLT-XPT",
			rawCode: "PLT",
			want:    "XPT",
		},
		{
			name:    "PLD-XPD",
			rawCode: "PLD",
			want:    "XPD",
		},
		{
			name:    "2_симв",
			rawCode: "PL",
			want:    "",
			wantErr: ErrWrongCurrencyCode,
		},
		{
			name:    "несуществующий",
			rawCode: "AAA",
			want:    "",
			wantErr: ErrNotExistingCurrency,
		},
		{
			name:    "DEM",
			rawCode: "DEM",
			want:    "DEM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ParseCurrencyCode(tt.rawCode)
			if tt.wantErr != nil {
				require.ErrorIs(t, tt.wantErr, gotErr)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
