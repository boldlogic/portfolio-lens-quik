package config

type Worker struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Name     string `yaml:"name" json:"name"`
	Interval uint16 `yaml:"interval_sec" json:"interval_sec"`
}

const defaultWorkerInterval uint16 = 60

func (w *Worker) applyDefaults(defaultName string) {
	if w.Name == "" {
		w.Name = defaultName
	}
	if w.Interval == 0 {
		w.Interval = defaultWorkerInterval
	}
}
