package service

type WorkerConfig struct {
	BatchSize uint16 `yaml:"batch_size" json:"batch_size"`
	QueueSize uint16 `yaml:"queue_size" json:"queue_size"`
	Interval  uint16 `yaml:"interval_ms" json:"interval_ms"`
}

const (
	defaultBatchSize uint16 = 100
	defaultQueueSize uint16 = 100
	defaultInterval  uint16 = 100
)

func (w *WorkerConfig) ApplyDefaults() {
	if w.BatchSize == 0 {
		w.BatchSize = defaultBatchSize
	}
	if w.QueueSize == 0 {
		w.QueueSize = defaultQueueSize
	}
	if w.Interval == 0 {
		w.Interval = defaultInterval
	}
}
