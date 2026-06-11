package application

import (
	"context"
	"sync"
	"time"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/periodic"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/config"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/features/dictionary"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/features/fxcbr"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/storage"
	storagemssql "github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/storage/mssql"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger
	store  *storage.Storage
	wg     sync.WaitGroup
}

const defaultConfigPath = "quik-currency/config.yaml"

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.Log)
	return &Application{
		cfg:    cfg,
		Logger: log,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	store, err := storage.NewStorage(ctx, a.cfg.Db, a.Logger)
	if err != nil {
		return err
	}
	a.store = store
	currencyRepo := storagemssql.NewCurrencyRepo(a.store.Db)
	curSvc := dictionary.NewService(a.Logger, currencyRepo)
	fxSvc := fxcbr.NewService(a.Logger, currencyRepo)

	workers := make([]periodic.Worker, 0, 2)

	if a.cfg.CurrencyWorkerConfig.Enabled {
		workers = append(workers, dictionary.NewUpdateCurrencyDictionaryWorker(
			curSvc,
			a.Logger,
			a.cfg.CurrencyWorkerConfig.Name,
			time.Duration(a.cfg.CurrencyWorkerConfig.Interval)*time.Second,
		))
	}

	if a.cfg.FxCBRWorkerConfig.Enabled {
		workers = append(workers, fxcbr.NewMergeFxCBRRatesQuikWorker(
			fxSvc,
			a.Logger,
			a.cfg.FxCBRWorkerConfig.Name,
			time.Duration(a.cfg.FxCBRWorkerConfig.Interval)*time.Second,
		))
	}

	if len(workers) == 0 {
		return nil
	}

	runner := periodic.NewRunner(workers...)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	return nil
}

func (a *Application) Wait(ctx context.Context) {
	<-ctx.Done()
	a.wg.Wait()
}
