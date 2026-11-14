package compose

import (
	"context"
	"fmt"
	"sync"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/config"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/jobs"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/persist"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/pipeline"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/schedule"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/sqlserver"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"go.uber.org/zap"
)

const defaultConfigPath = "configs/quik-instruments-worker-config.yaml"

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *sqlserver.Repository
	wg     sync.WaitGroup
}

func New() *App {
	return &App{}
}

func (a *App) Init() (*zap.Logger, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	a.cfg = cfg
	a.logger = logger.New(cfg.Log)
	return a.logger, nil
}

func (a *App) Start(ctx context.Context) error {
	if a.cfg == nil {
		return fmt.Errorf("приложение не проинициализировано")
	}

	db, err := sqlserver.NewRepository(ctx, a.cfg.Db, a.logger)
	if err != nil {
		return err
	}
	a.db = db

	instruments := persist.NewInstruments(db)

	runner := pipeline.NewRunner(
		a.logger,
		jobs.NewCurrencies(db, db, instruments),
		jobs.NewQuotes(db, db, instruments, quik.InstrumentTitleSecurities, "нормализация акций завершена"),
		jobs.NewQuotes(db, db, instruments, quik.InstrumentTitleBonds, "нормализация облигаций завершена"),
		jobs.NewQuotes(db, db, instruments, quik.InstrumentTitleCurrencyPairs, "нормализация валютных пар завершена"),
	)

	ticker := schedule.NewTicker(a.cfg.Worker, runner, a.logger)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		ticker.Run(ctx)
	}()

	return nil
}

func (a *App) Wait(ctx context.Context) {
	<-ctx.Done()
	a.wg.Wait()
	if a.db != nil {
		a.db.Close()
	}
}
