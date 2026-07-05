package debug

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.PprofConfig) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Addr(),
			Handler: mux,
		},
	}
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

func (s *Server) ListenAndServe() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("pprof server остановлен с ошибкой: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
