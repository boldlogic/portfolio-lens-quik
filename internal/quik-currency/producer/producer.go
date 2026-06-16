package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Producer struct {
	enabled      bool
	client       *kgo.Client
	brokers      []string
	topic        string
	flushTimeout time.Duration
	logger       *zap.Logger
}

func NewProducer(ctx context.Context, cfg Config, logger *zap.Logger) (*Producer, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	producer := Producer{
		enabled: cfg.Enabled,
		brokers: cfg.Brokers,
		topic:   cfg.Topic,
	}

	if !producer.enabled {
		return &producer, nil
	}

	opts := createOpts(cfg)
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("некорректный конфиг, %w", err)
	}
	err = client.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("брокер недоступен, %w", err)
	}

	producer.client = client
	producer.flushTimeout = time.Duration(cfg.Params.ProduceTimeout) * time.Second

	return &producer, nil
}

func createOpts(config Config) []kgo.Opt {
	opts := []kgo.Opt{
		kgo.SeedBrokers(config.Brokers...),
		kgo.ClientID(config.ClientID),
		kgo.DefaultProduceTopic(config.Topic),
	}
	switch config.Params.Acks {
	case "0":
		opts = append(opts, kgo.RequiredAcks(kgo.NoAck()))
	case "1":
		opts = append(opts, kgo.RequiredAcks(kgo.LeaderAck()))
	case "all":
		opts = append(opts, kgo.RequiredAcks(kgo.AllISRAcks()))
	}
	if config.Params.MaxRecordBatchSize != 0 {
		opts = append(opts, kgo.ProducerBatchMaxBytes(config.Params.MaxRecordBatchSize))
	}
	if config.Params.MaxBufferedRecords != 0 {
		opts = append(opts, kgo.MaxBufferedRecords(config.Params.MaxBufferedRecords))
	}
	if config.Params.MaxProduceOnFlight != 0 {
		opts = append(opts, kgo.MaxProduceRequestsInflightPerBroker(config.Params.MaxProduceOnFlight))
	}
	if config.Params.ProduceTimeout != 0 {
		opts = append(opts, kgo.ProduceRequestTimeout(time.Duration(config.Params.ProduceTimeout)*time.Second))
	}

	return opts
}

func (p *Producer) Close() {
	if !p.enabled {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.flushTimeout)
	defer cancel()

	if err := p.client.Flush(ctx); err != nil {
		p.logger.Error("не удалось дослать события в kafka при остановке", zap.Error(err))
	}
	p.client.Close()
}

func (p *Producer) Produce(ctx context.Context, records ...*kgo.Record) error {

	if err := p.client.ProduceSync(ctx, records...).FirstErr(); err != nil {
		return fmt.Errorf("не удалось записать сообщение: %w", err)
	}

	return nil
}

func (p *Producer) PublishCurrencies(ctx context.Context, currencies []quik.Currency) error {

	if !p.enabled {
		return nil
	}

	var records []*kgo.Record

	for _, currency := range currencies {
		event := currencyToEvent(currency)
		msg, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("не удалось закодировать сообщение")
		}
		record := kgo.Record{
			Topic: p.topic,
			Key:   []byte(event.ISOCharCode),
			Value: msg,
		}
		records = append(records, &record)
	}
	err := p.Produce(ctx, records...)
	if err != nil {
		p.logger.Error("не удалось отправить события",
			zap.Error(err),
		)
	}
	return nil
}
