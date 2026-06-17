package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boldlogic/packages/transport/httpserver/request"
	"github.com/boldlogic/packages/utils/dates"
	"github.com/stretchr/testify/assert"
)

func Test_parseClientCodesQueryParam(t *testing.T) {
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
		{
			name: "clientCodes_с_пробелами_и_пустыми_частями",
			r:    httptest.NewRequest(http.MethodGet, "/?clientCodes=AA,%20BB,%20,%20CC%20", nil),
			want: []string{"AA", "BB", "CC"},
		},
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			got := parseClientCodesQueryParam(testCase.r)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func Test_parseTargetCurrencyQueryParam(t *testing.T) {
	tests := []struct {
		name string
		r    *http.Request
		want *string
	}{
		{
			name: "не_передан",
			r:    httptest.NewRequest(http.MethodGet, "/", nil),
			want: nil,
		},
		{
			name: "пустой",
			r:    httptest.NewRequest(http.MethodGet, "/?targetCurrency=", nil),
			want: nil,
		},
		{
			name: "usd_нормализуется_в_upper",
			r:    httptest.NewRequest(http.MethodGet, "/?targetCurrency=usd", nil),
			want: strPtr("USD"),
		},
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			got := parseTargetCurrencyQueryParam(testCase.r)
			if testCase.want == nil {
				assert.Nil(t, got)
				return
			}
			assert.NotNil(t, got)
			assert.Equal(t, *testCase.want, *got)
		})
	}
}

func Test_parsePortfolioQueryParams(t *testing.T) {
	tests := []struct {
		name string
		r    *http.Request
		want portfolioQueryParams
	}{
		{
			name: "без_параметров",
			r:    httptest.NewRequest(http.MethodGet, "/", nil),
			want: portfolioQueryParams{},
		},
		{
			name: "targetCurrency_и_clientCodes",
			r:    httptest.NewRequest(http.MethodGet, "/?targetCurrency=EUR&clientCodes=AA,BB", nil),
			want: portfolioQueryParams{
				TargetCurrency: strPtr("EUR"),
				ClientCodes:    []string{"AA", "BB"},
			},
		},
		{
			name: "currency_не_маппится_на_targetCurrency",
			r:    httptest.NewRequest(http.MethodGet, "/?currency=USD", nil),
			want: portfolioQueryParams{},
		},
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			got := parsePortfolioQueryParams(testCase.r)
			if testCase.want.TargetCurrency == nil {
				assert.Nil(t, got.TargetCurrency)
			} else {
				assert.NotNil(t, got.TargetCurrency)
				assert.Equal(t, *testCase.want.TargetCurrency, *got.TargetCurrency)
			}
			assert.Equal(t, testCase.want.ClientCodes, got.ClientCodes)
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func Test_parseLimitsQueryParams(t *testing.T) {
	wantDate := time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name    string
		target  string
		want    limitsQueryParams
		wantErr error
	}{
		{
			name:   "все_параметры_переданы",
			target: "/?loadDate=2026-05-31&limit=25&offset=7&clientCodes=AA,%20BB&includeTotalCount=true",
			want: limitsQueryParams{
				LoadDate:          wantDate,
				Limit:             25,
				Offset:            7,
				ClientCodes:       []string{"AA", "BB"},
				IncludeTotalCount: true,
			},
		},
		{
			name:    "некорректный_loadDate",
			target:  "/?loadDate=2026-05-32",
			wantErr: dates.ErrWrongDateFormat,
		},
		{
			name:    "некорректный_limit",
			target:  "/?limit=0",
			wantErr: request.ErrInvalidLimit,
		},
		{
			name:    "некорректный_offset",
			target:  "/?offset=-1",
			wantErr: request.ErrInvalidOffset,
		},
		{
			name:    "некорректный_includeTotalCount",
			target:  "/?includeTotalCount=maybe",
			wantErr: request.ErrInvalidQuery,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseLimitsQueryParams(httptest.NewRequest(http.MethodGet, testCase.target, nil))

			if testCase.wantErr != nil {
				assert.True(t, errors.Is(err, testCase.wantErr), "err = %v", err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}
