package application

import (
	"context"
	"sync"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/config"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/repository"
	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/service"
	"go.uber.org/zap"
)

type Application struct {
	cfg     *config.Config
	logger  *zap.Logger
	svc     *service.Service
	repo    *repository.Repository
	server  *httpserver.Server
	errChan chan error
	wg      sync.WaitGroup
}

const (
	defaultConfigPath = "quik-currency/internal/configs/config.yaml"
	errChanBufSize    = 1
)

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return &Application{}, err
	}
	log := logger.New(cfg.Log)
	return &Application{
		cfg:     cfg,
		logger:  log,
		errChan: make(chan error, errChanBufSize),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.logger)
	if err != nil {
		return err
	}
	a.repo = repo

	a.svc = service.NewService(a.repo, a.logger)

	if err = a.svc.InitCurrencyDictionary(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Application) Wait(ctx context.Context, cancel context.CancelFunc) error {
	var appErr error

	errWg := sync.WaitGroup{}
	errWg.Add(1)

	go func() {
		defer errWg.Done()
		for err := range a.errChan {
			cancel()
			appErr = err
		}
	}()

	<-ctx.Done()

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
