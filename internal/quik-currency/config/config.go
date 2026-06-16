package config

import (
	"errors"
	"fmt"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/transport/httpserver"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/producer"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/worker"
)

type Config struct {
	Log               logger.Config           `yaml:"log" json:"log"`
	Db                dbconfig.DBConfig       `yaml:"db" json:"db"`
	Server            httpserver.ServerConfig `yaml:"server" json:"server"`
	FxCBRJobConfig    worker.JobConfig        `yaml:"fx_cbr_worker" json:"fx_cbr_worker"`
	CurrencyJobConfig worker.JobConfig        `yaml:"currency_worker" json:"currency_worker"`
	Kafka             producer.Config         `yaml:"kafka_producer" json:"kafka_producer"`
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
	kErrs := c.Kafka.Validate()
	if len(kErrs) > 0 {
		errs = append(errs, kErrs...)
	}
	return errs
}

func (c *Config) applyDefaults() {
	c.Db.ApplyDefaults()
	c.Db.ApplySecretsFromEnv()
	c.Server.ApplyDefaults()
	c.FxCBRJobConfig.ApplyDefaults("quik.currency.crossrates.import")
	c.CurrencyJobConfig.ApplyDefaults("quik.currency.dictionary.refresh")
}
