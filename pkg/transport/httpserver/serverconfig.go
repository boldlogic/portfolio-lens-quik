package httpserver

import (
	"fmt"
	"strings"
)

type ServerConfig struct {
	ListenHost   string           `yaml:"listen_host" json:"listen_host"`
	ExternalHost string           `yaml:"external_host" json:"external_host"`
	Port         int              `yaml:"port" json:"port"`
	Timeout      int              `yaml:"timeout" json:"timeout"`
	TLS          *TLSListenConfig `yaml:"tls,omitempty" json:"tls,omitempty"`
}

// TLSListenConfig включает HTTPS и взаимную проверку клиентских сертификатов (mTLS).
type TLSListenConfig struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	CertFile     string `yaml:"cert_file" json:"cert_file"`
	KeyFile      string `yaml:"key_file" json:"key_file"`
	ClientCAFile string `yaml:"client_ca_file" json:"client_ca_file"`
}

func (srv *ServerConfig) ApplyDefaults() {
	if srv.ListenHost == "" {
		srv.ListenHost = "127.0.0.1"
	}
	if srv.ExternalHost == "" {
		srv.ExternalHost = "localhost"
	}
	if srv.Port == 0 {
		srv.Port = 80
	}
	if srv.Timeout == 0 {
		srv.Timeout = 60
	}

}

func (srv *ServerConfig) Validate() []error {
	var errs []error
	if srv.Port < 1 || srv.Port > 65535 {
		errs = append(errs, fmt.Errorf("в блоке 'server' некорректный 'port': должен быть в диапазоне 1-65535, получено %d", srv.Port))
	}
	if srv.TLS != nil && srv.TLS.Enabled {
		if strings.TrimSpace(srv.TLS.CertFile) == "" {
			errs = append(errs, fmt.Errorf("в блоке 'server.tls' при enabled=true нужен 'cert_file'"))
		}
		if strings.TrimSpace(srv.TLS.KeyFile) == "" {
			errs = append(errs, fmt.Errorf("в блоке 'server.tls' при enabled=true нужен 'key_file'"))
		}
		if strings.TrimSpace(srv.TLS.ClientCAFile) == "" {
			errs = append(errs, fmt.Errorf("в блоке 'server.tls' при enabled=true нужен 'client_ca_file' (PEM доверенного CA для клиентов)"))
		}
	}
	return errs
}
