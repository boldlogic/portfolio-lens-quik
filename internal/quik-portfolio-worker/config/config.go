package config

import (
	"errors"
	"fmt"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
)

type Config struct {
	Log logger.Config     `yaml:"log" json:"log"`
	Db  dbconfig.DBConfig `yaml:"db" json:"db"`
}

func LoadConfig(configPath string) (*Config, error) {

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

//

func (c *Config) validate() []error {
	var errs []error

	dbErrs := c.Db.Validate()
	if len(dbErrs) > 0 {
		errs = append(errs, dbErrs...)
	}

	return errs
}

func (c *Config) applyDefaults() {
	c.Db.ApplyDefaults()
	c.Db.ApplySecretsFromEnv()

}
