package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbzap"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpclient"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver"
)

type Config struct {
	Log             logger.Config           `yaml:"log" json:"log"`
	Server          httpserver.ServerConfig `yaml:"server" json:"server"`
	Grpc            GrpcConfig              `yaml:"grpc" json:"grpc"`
	Db              dbzap.DBConfig          `yaml:"db" json:"db"`
	ServiceRegistry ServiceRegistryConfig   `yaml:"service_registry" json:"service_registry"`
}

type ServiceRegistryConfig struct {
	ManagerBaseURL       string                      `yaml:"manager_base_url" json:"manager_base_url"`
	HeartbeatIntervalSec int                         `yaml:"heartbeat_interval_sec" json:"heartbeat_interval_sec"`
	InstanceID           string                      `yaml:"instance_id" json:"instance_id"`
	GrpcPublicAddr       string                      `yaml:"grpc_public_addr" json:"grpc_public_addr"`
	APISecret            string                      `yaml:"api_secret" json:"api_secret"`
	Mtls                 httpclient.MtlsClientConfig `yaml:"mtls,omitempty" json:"mtls,omitempty"`
}

type GrpcConfig struct {
	Port int `yaml:"port" json:"port"`
}

func (g *GrpcConfig) ApplyDefaults() {
	if g.Port == 0 {
		g.Port = 5051
	}
}

func (g *GrpcConfig) Addr() string {
	return fmt.Sprintf(":%d", g.Port)
}

func Load(configPath string) (*Config, error) {

	cfg, err := commonconfig.DecodeConfigStrict[Config](configPath)

	if err != nil {
		return nil, err
	}
	cfg.applyDefaults()
	errs := cfg.validate()
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("некорректный конфиг: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() []error {
	var errs []error

	dbErrs := c.Db.Validate()
	if len(dbErrs) > 0 {
		errs = append(errs, dbErrs...)
	}
	srvErrs := c.Server.Validate()
	if len(srvErrs) > 0 {
		errs = append(errs, srvErrs...)
	}
	if c.ServiceRegistry.Mtls.Any() && !c.ServiceRegistry.Mtls.Enabled() {
		errs = append(errs, fmt.Errorf("service_registry.mtls: задайте ca_file, cert_file и key_file"))
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(c.ServiceRegistry.ManagerBaseURL)), "https://") && !c.ServiceRegistry.Mtls.Enabled() {
		errs = append(errs, fmt.Errorf("service_registry: при https manager_base_url нужен полный блок mtls"))
	}
	return errs
}

func (c *Config) applyDefaults() {
	c.Db.ApplyDefaults()
	c.Db.ApplySecretsFromEnv()
	c.Server.ApplyDefaults()
	c.Grpc.ApplyDefaults()
	if c.ServiceRegistry.HeartbeatIntervalSec == 0 {
		c.ServiceRegistry.HeartbeatIntervalSec = 10
	}
	if v := os.Getenv("SERVICE_REGISTRY_API_SECRET"); v != "" {
		c.ServiceRegistry.APISecret = v
	}
}
