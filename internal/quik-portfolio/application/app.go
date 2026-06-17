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
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/config"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/observability"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/repository"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/service"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/transport/grpc"
	grpcv1 "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/transport/grpc/v1"
	portfolioserver "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/transport/http"
	v1 "github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio/transport/http/v1"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	defaultConfigPath = "configs/quik-portfolio-config.yaml"
)

type Application struct {
	cfg    *config.Config //+
	Logger *zap.Logger

	server  *httpserver.Server
	grpcSrv *grpc.Server
	repo    *repository.Repository
	svc     *service.Service

	runGroup *errgroup.Group
	runCtx   context.Context

	shutdownOnce sync.Once
	shutdownErr  error
}

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
	reg := metrics.New()
	recorder := observability.New(reg)

	repo, err := repository.NewRepository(ctx, a.cfg.Db, a.Logger, recorder)
	if err != nil {
		return err
	}
	a.repo = repo
	if err := observability.RegisterDBStats(reg, a.repo.Db); err != nil {
		return err
	}

	a.svc = service.NewService(a.repo, a.Logger)

	commonHandler := handler.NewHandler()
	handler := v1.NewHandler(commonHandler, a.svc, a.Logger)
	r := portfolioserver.NewRouter(handler, a.Logger, reg)
	a.server = httpserver.NewServer(r, a.cfg.Server)

	grpcHandler := grpcv1.NewHandler(a.svc, a.Logger)
	grpcSrv, err := grpc.NewServer(a.cfg.Grpc.Addr(), grpcHandler, a.Logger, reg)
	if err != nil {
		return fmt.Errorf("ошибка создания gRPC server: %w", err)
	}
	a.grpcSrv = grpcSrv

	a.runGroup, a.runCtx = errgroup.WithContext(ctx)

	a.runGroup.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server остановлен с ошибкой: %w", err)
		}
		return nil
	})

	a.runGroup.Go(func() error {
		if err := a.grpcSrv.Start(); err != nil {
			return fmt.Errorf("gRPC server остановлен с ошибкой: %w", err)
		}
		return nil
	})

	return nil
}

func (a *Application) Wait() error {

	if a.runGroup == nil {
		return errors.New("приложение не запущено")
	}
	a.Logger.Debug("ждем")

	<-a.runCtx.Done()

	a.Logger.Debug("получили сигнал")

	a.shutdownOnce.Do(func() {
		a.shutdownErr = a.shutdown(context.Background()) ///
	})
	a.Logger.Debug("вызвали shutdownOnce")

	runErr := a.runGroup.Wait()
	a.Logger.Debug("дождались wait")

	var closeErr error
	if a.repo != nil {
		if cerr := a.repo.Close(); cerr != nil {
			closeErr = fmt.Errorf("ошибка закрытия БД: %w", cerr)
		}
	}
	a.Logger.Debug("завершили БД")

	return errors.Join(runErr, a.shutdownErr, closeErr)
}

func (a *Application) shutdown(ctx context.Context) error {

	var errs []error

	if a.server != nil {
		httpShutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Server.Opts.ShutdownTimeout)*time.Second)
		if err := a.server.Shutdown(httpShutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("ошибка остановки http: %w", err))
		}
		cancel()
	}
	a.Logger.Debug("завершили http")

	if a.grpcSrv != nil {
		grpcShutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Grpc.ShutdownTimeout)*time.Second)
		if err := a.grpcSrv.StopWithTimeout(grpcShutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("ошибка остановки grpc: %w", err))
		}
		cancel()
	}

	a.Logger.Debug("завершили grpc")

	return errors.Join(errs...)
}
