package config

import "fmt"

type PprofConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	ListenHost string `yaml:"listen_host" json:"listen_host"`
	Port       int    `yaml:"port" json:"port"`
}

func (c *PprofConfig) ApplyDefaults() {
	if c.ListenHost == "" {
		c.ListenHost = "127.0.0.1"
	}
	if c.Port == 0 {
		c.Port = 6060
	}
}

func (c *PprofConfig) Validate() []error {
	if !c.Enabled {
		return nil
	}
	if c.Port < 1 || c.Port > 65535 {
		return []error{fmt.Errorf("в блоке 'pprof' некорректный 'port': должен быть в диапазоне 1-65535, получено %d", c.Port)}
	}
	return nil
}

func (c *PprofConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.ListenHost, c.Port)
}
