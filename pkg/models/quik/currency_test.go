package quik

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyFromQuik(t *testing.T) {
	tests := []struct {
		testName     string
		charCode     string
		name         *string
		wantCharCode string
		wantName     *string
		wantAlpha    CurrencyCode
		wantNumeric  int16
		wantErr      error
	}{
		{
			testName:     "обычная_валюта",
			charCode:     "CNY",
			name:         new("Юань"),
			wantCharCode: "CNY",
			wantName:     new("Юань"),
			wantAlpha:    CurrencyCode("CNY"),
			wantNumeric:  int16(156),
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, gotErr := CurrencyFromQuik(tt.charCode, tt.name)
			if tt.wantErr != nil {
				require.ErrorIs(t, tt.wantErr, gotErr)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, tt.wantCharCode, got.charCode)
			assert.Equal(t, tt.wantName, got.name)
			assert.Equal(t, tt.wantAlpha, got.alpha)
			assert.Equal(t, tt.wantNumeric, got.numeric)
		})
	}
}

func Test_currencyIso_Alpha(t *testing.T) {
	tests := []struct {
		name     string
		charCode string
		want     CurrencyCode
	}{
		{
			name:     "юань",
			charCode: "CNY",
			want:     CurrencyCode("CNY"),
		},
		{
			name:     "SUR",
			charCode: "SUR",
			want:     CurrencyCode("RUB"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := CurrencyFromQuik(tt.charCode, nil)
			require.NoError(t, gotErr)

			assert.Equal(t, tt.want, got.alpha)
			assert.Equal(t, got.alpha, got.Alpha())

			assert.NotEmpty(t, got.alpha.String())
		})
	}
}

func Test_currencyIso_MinorUnits(t *testing.T) {
	tests := []struct {
		name     string
		charCode string
		want     int32
	}{
		{
			name:     "SUR_2",
			charCode: "SUR",
			want:     int32(2),
		},
		{
			name:     "GLD-1",
			charCode: "GLD",
			want:     int32(0),
		},
		{
			name:     "ROL-1",
			charCode: "ROL",
			want:     int32(0),
		},
		{
			name:     "JPY_0",
			charCode: "JPY",
			want:     int32(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := CurrencyFromQuik(tt.charCode, nil)
			require.NoError(t, gotErr)

			assert.Equal(t, tt.want, got.minorUnits)
			assert.Equal(t, got.minorUnits, got.MinorUnits())
		})
	}
}
