package config

import (
	"errors"
	"fmt"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
)

type Config struct {
	Log    logger.Config     `yaml:"log" json:"log"`
	Db     dbconfig.DBConfig `yaml:"db" json:"db"`
	Worker WorkerConfig      `yaml:"worker" json:"worker"`
}

type WorkerConfig struct {
	Enabled       bool   `yaml:"enabled" json:"enabled"`
	Name          string `yaml:"name" json:"name"`
	Interval      uint16 `yaml:"interval_sec" json:"interval_sec"`
	RunOnStart    bool   `yaml:"run_on_start" json:"run_on_start"`
	Timeout       uint16 `yaml:"timeout_sec" json:"timeout_sec"`
	MaxErrorCount uint16 `yaml:"max_error_count" json:"max_error_count"`
	BatchSize     int    `yaml:"batch_size" json:"batch_size"`
}

const (
	defaultWorkerName          = "quik.instruments.normalize"
	defaultWorkerInterval      = 60
	defaultWorkerTimeout       = 30
	defaultWorkerMaxErrorCount = 10
	defaultBatchSize           = 500
)

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

func (c *Config) applyDefaults() {
	c.Db.ApplyDefaults()
	c.Db.ApplySecretsFromEnv()
	c.Worker.applyDefaults()
}

func (c *Config) validate() []error {
	var errs []error

	dbErrs := c.Db.Validate()
	if len(dbErrs) > 0 {
		errs = append(errs, dbErrs...)
	}
	if c.Worker.Enabled && c.Worker.BatchSize <= 0 {
		errs = append(errs, fmt.Errorf("worker.batch_size должен быть больше нуля"))
	}

	return errs
}

func (w *WorkerConfig) applyDefaults() {
	if !w.Enabled {
		return
	}
	if w.Name == "" {
		w.Name = defaultWorkerName
	}
	if w.Interval == 0 {
		w.Interval = defaultWorkerInterval
	}
	if w.Timeout == 0 {
		w.Timeout = defaultWorkerTimeout
	}
	if w.MaxErrorCount == 0 {
		w.MaxErrorCount = defaultWorkerMaxErrorCount
	}
	if w.BatchSize == 0 {
		w.BatchSize = defaultBatchSize
	}
}
