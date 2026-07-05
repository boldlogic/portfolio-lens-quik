package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"go.uber.org/zap"
)

type currencyRepo interface {
	SelectCurrencies(ctx context.Context) ([]models.Currency, error)
	SelectNewCurrenciesFromCrossrates(ctx context.Context) ([]models.CurrencyFromCrossrates, error)
	MergeCurrencies(ctx context.Context, currencies []quik.Currency) error
}

type CurrencyWorkerConfig struct {
	ExternalName string `yaml:"name" json:"name"`
	Interval     uint16 `yaml:"interval_ms" json:"interval_ms"`
}

const (
	defaultWorkerInterval uint16 = 60
)

func (w *CurrencyWorkerConfig) ApplyDefaults() {
	if w.ExternalName == "" {
		w.ExternalName = "quik.currency.dictionary.refresh"
	}
	if w.Interval == 0 {
		w.Interval = defaultWorkerInterval
	}
}

type currencyWorker struct {
	externalName string
	repo         currencyRepo
	interval     time.Duration
	cache        *currencyCache
	logger       *zap.Logger
}

func (w currencyWorker) run(ctx context.Context) error {

	err := w.warmup(ctx)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := w.refresh(ctx)
			if err != nil {
				w.logger.Error(err.Error())
			}
			w.logger.Debug("", zap.Any("", unsafe.Sizeof(w.cache.currencies)))
			w.logger.Debug("", zap.Any("", w.cache.currencies))
		case <-ctx.Done():
			return nil
		}
	}

}

func newCurrencyWorker(cfg CurrencyWorkerConfig, repo currencyRepo, logger *zap.Logger) *currencyWorker {
	return &currencyWorker{
		externalName: cfg.ExternalName,
		repo:         repo,
		interval:     time.Duration(cfg.Interval) * time.Millisecond,
		cache:        newCurrencyCache(),
		logger:       logger,
	}
}

type currencyCache struct {
	currencies map[quik.CurrencyCode]quik.Currency
	mu         sync.RWMutex
}

func newCurrencyCache() *currencyCache {
	return &currencyCache{
		currencies: make(map[quik.CurrencyCode]quik.Currency, 200),
	}
}

func (w currencyWorker) warmup(ctx context.Context) error {

	cs, err := w.repo.SelectCurrencies(ctx)
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}
	if len(cs) > 0 {
		for _, c := range cs {
			ccy, err := quik.CurrencyFromQuik(c.CharCode, c.Name)
			if err != nil {
				continue
			}

			w.cache.mu.Lock()
			w.cache.currencies[ccy.Alpha()] = ccy
			w.cache.mu.Unlock()
		}

		return nil
	}
	curs, err := w.getQuikCurrecies(ctx)
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}
	err = w.upsertCurrencies(ctx, curs)
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}
	for _, c := range curs {
		w.cache.mu.Lock()
		w.cache.currencies[c.Alpha()] = c
		w.cache.mu.Unlock()
	}

	return nil

}

func (w *currencyWorker) getQuikCurrecies(ctx context.Context) ([]quik.Currency, error) {
	raw, err := w.repo.SelectNewCurrenciesFromCrossrates(ctx)
	if err != nil {
		w.logger.Error(err.Error())
		return nil, fmt.Errorf("ошибка при получении кросс-курсов")
	}
	out := make([]quik.Currency, 0, len(raw))
	dedup := make(map[quik.CurrencyCode]struct{})

	for _, ccy := range raw {
		currency, err := quik.CurrencyFromQuik(ccy.IsoCharCode, &ccy.Name)
		if err != nil {
			continue
		}
		if _, ok := dedup[currency.Alpha()]; ok {
			continue
		}
		dedup[currency.Alpha()] = struct{}{}
		out = append(out, currency)
	}
	return out, nil
}

func (w *currencyWorker) upsertCurrencies(ctx context.Context, cur []quik.Currency) error {
	return w.repo.MergeCurrencies(ctx, cur)
}

func (w *currencyWorker) refresh(ctx context.Context) error {
	curs, err := w.getQuikCurrecies(ctx)
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}
	forUpd := make([]quik.Currency, 0, len(curs))
	for _, c := range curs {
		w.cache.mu.RLock()
		cur, ok := w.cache.currencies[c.Alpha()]
		w.cache.mu.RUnlock()
		if !ok {
			w.logger.Debug("новая валюта", zap.String("", c.Alpha().String()))
			forUpd = append(forUpd, cur)
		}
	}
	if len(forUpd) == 0 {
		return nil
	}

	err = w.upsertCurrencies(ctx, forUpd)
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}
	for _, c := range forUpd {
		w.cache.mu.Lock()
		w.cache.currencies[c.Alpha()] = c
		w.cache.mu.Unlock()

	}
	return nil
}
