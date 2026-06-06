package application

import (
	"context"
	"sync"
	"time"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/periodic"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/config"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/repository"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/service"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/workers"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	logger *zap.Logger
	svc    *service.Service
	repo   *repository.Repository
	wg     sync.WaitGroup
}

const defaultConfigPath = "quik-currency/internal/configs/config.yaml"

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return &Application{}, err
	}
	log := logger.New(cfg.Log)
	return &Application{
		cfg:    cfg,
		logger: log,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	repo, err := repository.NewRepository(ctx, a.cfg.Db, a.logger)
	if err != nil {
		return err
	}
	a.repo = repo

	a.svc = service.NewService(a.repo, a.logger)

	if err = a.svc.InitCurrencyDictionary(ctx); err != nil {
		return err
	}

	runner := periodic.NewRunner(
		workers.NewMergeFxCBRRatesQuikWorker(a.svc, a.logger, 60*time.Second),
	)
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
