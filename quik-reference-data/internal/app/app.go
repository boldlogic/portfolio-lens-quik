package app

import (
	"context"
	"sync"
	"time"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/periodic"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/entities/firm"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/features/syncfirms"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/shared/config"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/shared/db"
	"go.uber.org/zap"
)

const (
	defaultConfigPath = "quik-reference-data/config.yaml"
	errChanBufSize    = 1
)

type application struct {
	config  *config.Config
	logger  *zap.Logger
	repo    *db.Repository
	errChan chan error
	wg      sync.WaitGroup
}

func New() *application {
	return &application{
		errChan: make(chan error, errChanBufSize),
	}
}

func (a *application) Init() (*zap.Logger, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)

	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return &zap.Logger{}, err
	}
	a.config = conf
	a.logger = logger.New(conf.Log)

	return a.logger, nil
}

func (a *application) Start(ctx context.Context) error {

	repo, err := db.NewRepository(ctx, a.config.Db.GetDSN(), a.logger)
	if err != nil {
		return err
	}
	a.repo = repo

	return nil
}

func (a *application) Wait(ctx context.Context, cancel context.CancelFunc) error {
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

func (a *application) InitFeatures(ctx context.Context) error {
	firmRepo := firm.NewFirmsRepo(a.repo)

	firmSvc := syncfirms.NewService(firmRepo, a.logger)
	runner := periodic.NewRunner(
		syncfirms.NewActualizeFirmsWorker(firmSvc, a.logger, 60*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	return nil
}
