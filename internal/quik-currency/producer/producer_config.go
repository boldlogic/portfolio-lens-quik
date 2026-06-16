package producer

import (
	"fmt"
	"slices"
)

type Config struct {
	Enabled  bool     `yaml:"enabled" json:"enabled"`
	Brokers  []string `yaml:"brokers" json:"brokers"`
	Topic    string   `yaml:"topic" json:"topic"`
	ClientID string   `yaml:"client_id" json:"client_id"`
	Params   Params   `yaml:"params" json:"params"`
}

type Params struct {
	Acks               acks   `yaml:"acks" json:"acks,omitempty"`
	MaxRecordBatchSize int32  `yaml:"max_record_batch_size" json:"max_record_batch_size,omitempty"`
	MaxBufferedRecords int    `yaml:"max_buffered_records" json:"max_buffered_records,omitempty"`
	MaxProduceOnFlight int    `yaml:"max_produce_on_flight" json:"max_produce_on_flight,omitempty"`
	ProduceTimeout     uint16 `yaml:"produce_timeout" json:"produce_timeout,omitempty"`
}

func (p *Params) Validate() []error {

	var errs []error

	if err := p.Acks.validate(); err != nil {
		errs = append(errs, fmt.Errorf("недопустимое значение acks"))
	}

	if p.MaxRecordBatchSize < 0 {
		errs = append(errs, fmt.Errorf("max_record_batch_size не может быть меньше нуля"))
	}
	if p.MaxBufferedRecords < 0 {
		errs = append(errs, fmt.Errorf("max_buffered_records не может быть меньше нуля"))
	}
	if p.MaxProduceOnFlight < 0 {
		errs = append(errs, fmt.Errorf("max_produce_on_flight не может быть меньше нуля"))
	}

	return errs

}

type acks string

func (a acks) validate() error {
	switch a {
	case "0":
		return nil
	case "1":
		return nil
	case "all":
		return nil
	default:
		return fmt.Errorf("")
	}

}

func (c *Config) Validate() []error {
	if !c.Enabled {
		return nil
	}
	var errs []error
	if len(c.Brokers) == 0 {
		errs = append(errs, fmt.Errorf("список brokers не может быть пустым"))
	}
	if slices.Contains(c.Brokers, "") {
		errs = append(errs, fmt.Errorf("адрес брокера не может быть пустой строкой"))
	}
	if c.Topic == "" {
		errs = append(errs, fmt.Errorf("topic не может быть пустым"))
	}
	if c.ClientID == "" {
		errs = append(errs, fmt.Errorf("client_id не может быть пустым"))
	}
	opErrs := c.Params.Validate()
	if len(opErrs) != 0 {
		errs = append(errs, opErrs...)
	}
	return errs

}
