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
	"github.com/boldlogic/packages/metrics"
	"github.com/boldlogic/packages/transport/httpserver"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/config"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/repository"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/service"
	writeserver "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/transport/http"
	v1 "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/transport/http/v1"
	"go.uber.org/zap"
)

const (
	defaultConfigPath = "quik-portfolio-writer/config.yaml"
	errChanBufSize    = 1
)

type application struct {
	config  *config.Config
	logger  *zap.Logger
	repo    *repository.Repository
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
	repo, err := repository.NewRepository(ctx, a.config.Db, a.logger)
	if err != nil {
		return err
	}
	a.repo = repo

	svc := service.NewService(a.repo, a.logger)
	reg := metrics.New()
	commonHandler := handler.NewHandler()
	v1 := v1.NewHandler(commonHandler, svc, a.logger)
	router := writeserver.NewRouter(v1, a.logger, reg)
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
