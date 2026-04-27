package middleware

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (m Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		headers := make([]string, 0, len(r.Header))
		for k, v := range r.Header {
			headers = append(headers, k+":"+strings.Join(v, ", "))
		}

		rw := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		m.logger.Debug("запрос",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.Int("status", rw.status),
			zap.String("headers", strings.Join(headers, ",")),
			zap.Duration("duration", time.Since(start)))
	})
}
