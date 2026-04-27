package httpserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	httpServer *http.Server
	addr       string
	tlsListen  *TLSListenConfig
}

func NewServer(handler http.Handler, cfg ServerConfig) *Server {
	listenHost := cfg.ListenHost
	if listenHost == "" {
		listenHost = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", listenHost, cfg.Port)
	timeout := time.Duration(cfg.Timeout) * time.Second

	srv := &Server{
		addr: addr,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       timeout,
			ReadHeaderTimeout: timeout,
			WriteTimeout:      timeout,
			IdleTimeout:       timeout,
		},
	}
	if cfg.TLS != nil && cfg.TLS.Enabled {
		srv.tlsListen = cfg.TLS
	}
	return srv
}

func (s *Server) ListenAndServe() error {
	if s.tlsListen != nil {
		tlsCfg, err := buildServerTLSConfig(s.tlsListen)
		if err != nil {
			return fmt.Errorf("tls listen: %w", err)
		}
		s.httpServer.TLSConfig = tlsCfg
		return s.httpServer.ListenAndServeTLS("", "")
	}
	return s.httpServer.ListenAndServe()
}

func buildServerTLSConfig(cfg *TLSListenConfig) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("загрузка server cert/key: %w", err)
	}
	pemData, err := os.ReadFile(cfg.ClientCAFile)
	if err != nil {
		return nil, fmt.Errorf("чтение client_ca_file: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pemData) {
		return nil, fmt.Errorf("client_ca_file: нет валидных PEM-сертификатов")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
