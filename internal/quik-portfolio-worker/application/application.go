package application

import (
	"context"
	"sync"
	"time"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/periodic"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-worker/config"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-worker/repository"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-worker/service"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-worker/workers"
	"go.uber.org/zap"
)

const defaultConfigPath = "configs/quik-portfolio-worker-config.yaml"

type application struct {
	config *config.Config
	logger *zap.Logger
	wg     sync.WaitGroup
}

func New() *application {
	return &application{}
}

func (a *application) Init() (*zap.Logger, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	a.config = conf
	a.logger = logger.New(conf.Log)
	return a.logger, nil
}

func (a *application) Start(ctx context.Context) error {
	repo, err := repository.NewRepository(ctx, a.config.Db, a.logger)
	if err != nil {
		return err
	}

	svc := service.NewService(repo, a.logger)
	runner := periodic.NewRunner(
		workers.NewRollForwardMoneyLimitsWorker(svc, a.logger, 60*time.Second),
		workers.NewRollForwardSecurityLimitsWorker(svc, a.logger, 60*time.Second),
		workers.NewRollForwardOtcWorker(svc, a.logger, 60*time.Second),
	)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	return nil
}

func (a *application) Wait(ctx context.Context) {
	<-ctx.Done()
	a.wg.Wait()
}
