package quik_test

import (
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/stretchr/testify/assert"
)

func TestParseCurrencyCode(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawCode string
		want    quik.CurrencyCode
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "обычная валюта",
			rawCode: "CNY",
			want:    "CNY",
		},
		{
			name:    "GLD",
			rawCode: "GLD",
			want:    "XAU",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := quik.ParseCurrencyCode(tt.rawCode)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseCurrencyCode() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseCurrencyCode() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			assert.Equal(t, tt.want, got)
		})
	}
}
