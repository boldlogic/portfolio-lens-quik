package config

import (
	"errors"
	"fmt"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbzap"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver"
)

type Config struct {
	Log    logger.Config           `yaml:"log" json:"log"`
	Server httpserver.ServerConfig `yaml:"server" json:"server"`
	Db     dbzap.DBConfig          `yaml:"db" json:"db"`
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
	return errs
}

func (c *Config) applyDefaults() {
	c.Db.ApplyDefaults()
	c.Db.ApplySecretsFromEnv()
	c.Server.ApplyDefaults()
}
