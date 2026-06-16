package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/transport/httpserver"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/config"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/features/marketdata"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/features/reference"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/producer"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/storage"
	storagemssql "github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/storage/mssql"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/transport/workerserver"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/worker"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger
	store  *storage.Storage

	repo     *storagemssql.CurrencyRepo
	server   *httpserver.Server
	producer *producer.Producer

	runGroup     *errgroup.Group
	runCtx       context.Context
	shutdownOnce sync.Once
	shutdownErr  error
}

const defaultConfigPath = "configs/quik-currency-config.yaml"

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

	a.repo = storagemssql.NewCurrencyRepo(a.store.Db)
	handler := workerserver.NewHandler(a.store.Db)
	router := workerserver.NewRouter(handler)
	a.server = httpserver.NewServer(router, a.cfg.Server)
	pr, err := producer.NewProducer(ctx, a.cfg.Kafka, a.Logger)
	if err != nil {
		return err
	}
	a.producer = pr

	runner := worker.NewRunner(a.prepareWorkers()...)

	a.runGroup, a.runCtx = errgroup.WithContext(ctx)
	a.runGroup.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server остановлен с ошибкой: %w", err)

		}
		return nil
	})
	a.runGroup.Go(func() error {
		runner.Run(a.runCtx)
		return nil
	})

	return nil
}

func (a *Application) Wait(ctx context.Context) error {
	if a.runGroup == nil {
		return fmt.Errorf("приложение не запущено")
	}
	<-a.runCtx.Done()
	a.shutdownOnce.Do(func() {
		a.shutdownErr = a.shutdown(context.Background())
	})
	return a.runGroup.Wait()
}

func (a *Application) prepareWorkers() []worker.Worker {
	currencyCatalog := reference.NewCurrencyCatalog(a.Logger, a.repo, a.producer)
	rateImporter := marketdata.NewRateImporter(a.Logger, a.repo)

	workers := make([]worker.Worker, 0, 2)
	currencyJob := worker.NewJob(a.cfg.CurrencyJobConfig, a.Logger, currencyCatalog.RefreshQuikCurrencies)
	fxJob := worker.NewJob(a.cfg.FxCBRJobConfig, a.Logger, rateImporter.ImportQuikCrossRates)
	workers = append(workers, currencyJob, fxJob)
	return workers
}

func (a *Application) shutdown(ctx context.Context) error {
	var errs []error
	if a.server != nil {
		shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Server.Opts.ShutdownTimeout)*time.Second)
		err := a.server.Shutdown(shCtx)
		if err != nil {
			errs = append(errs, fmt.Errorf("ошибка остановки http: %w", err))
		}
		cancel()
	}
	if a.producer != nil {
		a.producer.Close()
	}
	if a.store != nil {
		err := a.store.Db.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("ошибка закрытия БД: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
