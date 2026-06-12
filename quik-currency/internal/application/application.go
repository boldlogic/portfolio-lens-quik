package application

import (
	"context"
	"sync"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/config"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/features/marketdata"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/features/reference"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/storage"
	storagemssql "github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/storage/mssql"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/worker"
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

	repo := storagemssql.NewCurrencyRepo(a.store.Db)

	currencyCatalog := reference.NewCurrencyCatalog(a.Logger, repo)
	rateImporter := marketdata.NewRateImporter(a.Logger, repo)

	workers := make([]worker.Worker, 0, 2)
	currencyJob := worker.NewJob(a.cfg.CurrencyJobConfig, a.Logger, currencyCatalog.RefreshQuikCurrencies)
	fxJob := worker.NewJob(a.cfg.FxCBRJobConfig, a.Logger, rateImporter.ImportQuikCrossRates)
	workers = append(workers, currencyJob, fxJob)

	runner := worker.NewRunner(workers...)
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
