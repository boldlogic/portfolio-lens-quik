package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/metrics"
	"github.com/boldlogic/packages/periodic"
	"github.com/boldlogic/packages/transport/httpserver"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	referencehttp "github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/api/http"
	apiv1 "github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/api/http/v1"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/entities/firm"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/features/readfirms"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/features/syncfirms"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/features/writefirms"
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
	server  *httpserver.Server
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

	firmRepo := firm.NewFirmsRepo(a.repo)

	readFirmsSvc := readfirms.NewService(firmRepo, a.logger)
	writeFirmsSvc := writefirms.NewService(firmRepo, a.logger)
	syncFirmsSvc := syncfirms.NewService(firmRepo, a.logger)

	runner := periodic.NewRunner(
		syncfirms.NewActualizeFirmsWorker(syncFirmsSvc, a.logger, 60*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	reg := metrics.New()
	commonHandler := handler.NewHandler()
	apiHandler := apiv1.NewHandler(commonHandler, readFirmsSvc, writeFirmsSvc, a.logger)
	router := referencehttp.NewRouter(apiHandler, a.logger, reg)
	a.server = httpserver.NewServer(router, a.config.Server)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- fmt.Errorf("http server остановлен с ошибкой: %w", err)
		}
	}()

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

	if a.server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Duration(a.config.Server.Opts.ShutdownTimeout)*time.Second)
		defer shutdownCancel()
		_ = a.server.Shutdown(shutdownCtx)
	}

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
