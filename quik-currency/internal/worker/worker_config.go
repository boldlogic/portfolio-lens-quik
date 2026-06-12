package worker

type JobConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	Name       string `yaml:"name" json:"name"`
	Interval   uint16 `yaml:"interval_sec" json:"interval_sec"`
	RunOnStart bool   `yaml:"run_on_start" json:"run_on_start"`
	Timeout    uint16 `yaml:"timeout_sec" json:"timeout_sec"`
}

func (w *JobConfig) ApplyDefaults(defaultName string) {
	if w.Enabled == false {
		return
	}
	if w.Name == "" {
		w.Name = defaultName
	}
	if w.Interval == 0 {
		w.Interval = defaultWorkerInterval
	}
	if w.Timeout == 0 {
		w.Timeout = defaultWorkerTimeout
	}
}

const (
	defaultWorkerInterval uint16 = 60
	defaultWorkerTimeout  uint16 = 10
)
