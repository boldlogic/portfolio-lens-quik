package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// MtlsClientConfig задаёт пути к PEM для исходящих HTTPS-запросов с mTLS.
type MtlsClientConfig struct {
	CAFile     string `yaml:"ca_file" json:"ca_file"`
	CertFile   string `yaml:"cert_file" json:"cert_file"`
	KeyFile    string `yaml:"key_file" json:"key_file"`
	TimeoutSec int    `yaml:"timeout_sec" json:"timeout_sec"`
}

// Enabled true, если заданы все три PEM-пути.
func (c MtlsClientConfig) Enabled() bool {
	return strings.TrimSpace(c.CAFile) != "" &&
		strings.TrimSpace(c.CertFile) != "" &&
		strings.TrimSpace(c.KeyFile) != ""
}

// Any true, если задано хотя бы одно поле (для проверки «неполной» конфигурации).
func (c MtlsClientConfig) Any() bool {
	return strings.TrimSpace(c.CAFile) != "" ||
		strings.TrimSpace(c.CertFile) != "" ||
		strings.TrimSpace(c.KeyFile) != ""
}

// OptionalMtlsHTTPClient возвращает клиент с mTLS, если конфиг полный; (nil, nil), если mTLS не задан;
// ошибку при частично заполненном блоке.
func OptionalMtlsHTTPClient(c MtlsClientConfig) (*http.Client, error) {
	if !c.Any() {
		return nil, nil
	}
	if !c.Enabled() {
		return nil, fmt.Errorf("mtls: задайте все поля ca_file, cert_file, key_file")
	}
	return NewMtlsHTTPClient(c)
}

// NewMtlsHTTPClient собирает *http.Client с проверкой сервера по CA и клиентским сертификатом.
func NewMtlsHTTPClient(c MtlsClientConfig) (*http.Client, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("mtls: неполная конфигурация")
	}
	caPEM, err := os.ReadFile(c.CAFile)
	if err != nil {
		return nil, fmt.Errorf("mtls: чтение ca_file: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("mtls: ca_file без валидных PEM-сертификатов")
	}
	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("mtls: загрузка cert/key: %w", err)
	}
	tlsCfg := &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	tr := &http.Transport{TLSClientConfig: tlsCfg}
	timeout := 15 * time.Second
	if c.TimeoutSec > 0 {
		timeout = time.Duration(c.TimeoutSec) * time.Second
	}
	return &http.Client{Transport: tr, Timeout: timeout}, nil
}
