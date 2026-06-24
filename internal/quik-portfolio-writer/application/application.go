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
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/config"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/repository"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/service"
	writeserver "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/transport/http"
	v1 "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/transport/http/v1"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	defaultConfigPath = "configs/quik-portfolio-writer-config.yaml"
	errChanBufSize    = 2
)

type application struct {
	config *config.Config
	Logger *zap.Logger
	repo   *repository.Repository
	server *httpserver.Server
	svc    *service.Service

	runGroup     *errgroup.Group
	runCtx       context.Context
	shutdownOnce sync.Once
	shutdownErr  error
}

func New() (*application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)

	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &application{
		config: conf,
		Logger: logger.New(conf.Log),
	}, nil
}

func (a *application) Start(ctx context.Context) error {
	repo, err := repository.NewRepository(ctx, a.config.Db, a.Logger)
	if err != nil {
		return err
	}
	a.repo = repo

	a.svc = service.NewService(a.repo, a.config.Worker, a.Logger)
	reg := metrics.New()
	commonHandler := handler.NewHandler()
	v1 := v1.NewHandler(commonHandler, a.svc, a.Logger)
	router := writeserver.NewRouter(v1, a.Logger, reg)
	a.server = httpserver.NewServer(router, a.config.Server)

	a.runGroup, a.runCtx = errgroup.WithContext(ctx)

	a.runGroup.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server остановлен с ошибкой: %w", err)
		}
		return nil
	})

	a.runGroup.Go(func() error {
		if err := a.svc.Run(ctx); err != nil {
			return fmt.Errorf("воркер остановлен с ошибкой: %w", err)
		}
		return nil
	})
	return nil

}

func (a *application) Wait(ctx context.Context) error {
	var errs []error
	if a.runGroup == nil {
		return fmt.Errorf("приложение не запущено")
	}
	<-a.runCtx.Done()
	var err error
	a.shutdownOnce.Do(func() {
		err = a.shutdown(ctx)
	})
	if err != nil {
		errs = append(errs, err)
	}
	err = a.runGroup.Wait()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		errors.Join(errs...)
	}
	return nil
}

func (a *application) shutdown(ctx context.Context) error {
	var errs []error
	if a.server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Server.Opts.ShutdownTimeout)*time.Second)
		err := a.server.Shutdown(shutdownCtx)
		if err != nil {
			errs = append(errs, fmt.Errorf("ошибка остановки http: %w", err))
		}
		cancel()
	}
	if a.svc != nil {
		a.svc.Stop()
	}
	if a.repo != nil {
		err := a.repo.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("ошибка закрытия БД: %w", err))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil

}
