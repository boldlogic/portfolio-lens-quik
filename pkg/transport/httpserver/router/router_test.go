package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRouter_Health(t *testing.T) {
	r := NewRouter(zap.NewNop(), prometheus.NewRegistry())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, `"ok"`, rec.Body.String())
	contentTypes := rec.Result().Header.Values("Content-Type")
	require.NotEmpty(t, contentTypes)
	require.Contains(t, contentTypes[0], "application/json")
	require.Contains(t, contentTypes[0], "charset=UTF-8")
}
