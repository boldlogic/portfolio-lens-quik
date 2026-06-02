package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractClientsQueryParam(t *testing.T) {
	tests := []struct {
		name string
		r    *http.Request
		want []string
	}{
		{
			name: "clientCodes_передан",
			r:    httptest.NewRequest(http.MethodGet, "/?clientCodes=AA", nil),
			want: []string{"AA"},
		},
		{
			name: "не_передан",
			r:    httptest.NewRequest(http.MethodGet, "/", nil),
			want: nil,
		},
		{
			name: "пустой",
			r:    httptest.NewRequest(http.MethodGet, "/?clientCodes=", nil),
			want: nil,
		},
		{
			name: "несколько_clientCodes",
			r:    httptest.NewRequest(http.MethodGet, "/?clientCodes=AA,BB,CC", nil),
			want: []string{"AA", "BB", "CC"},
		},
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			got := extractClientsQueryParam(testCase.r)
			assert.Equal(t, testCase.want, got)
		})
	}
}
