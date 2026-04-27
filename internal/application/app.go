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
	"github.com/boldlogic/packages/periodic"
	"github.com/boldlogic/quik-portfolio/internal/config"
	"github.com/boldlogic/quik-portfolio/internal/repository"
	"github.com/boldlogic/quik-portfolio/internal/service"
	"github.com/boldlogic/quik-portfolio/internal/transport/grpc"
	grpcv1 "github.com/boldlogic/quik-portfolio/internal/transport/grpc/v1"
	portfolioserver "github.com/boldlogic/quik-portfolio/internal/transport/http"
	v1 "github.com/boldlogic/quik-portfolio/internal/transport/http/v1"
	"github.com/boldlogic/quik-portfolio/internal/workers"
	"github.com/boldlogic/quik-portfolio/pkg/registryinstance"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpclient"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

const (
	defaultConfigPath = "internal/configs/config.yaml"
	errChanBufSize    = 1
)

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger

	svc *service.Service

	errChan chan error
	wg      sync.WaitGroup
	repo    *repository.Repository

	server  *httpserver.Server
	grpcSrv *grpc.Server
}

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.Log)
	return &Application{
		cfg:     cfg,
		Logger:  log,
		errChan: make(chan error, errChanBufSize),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {

	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.Logger)
	if err != nil {
		return err
	}
	a.repo = repo

	a.svc = service.NewService(a.repo, a.Logger)

	runner := periodic.NewRunner(
		workers.NewRollForwardMoneyLimitsWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewRollForwardSecurityLimitsWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewRollForwardOtcWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewActualizeFirmsWorker(a.svc, a.Logger, 60*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	reg := metrics.New()
	commonHandler := handler.NewHandler()
	handler := v1.NewHandler(commonHandler, a.svc, a.Logger)
	r := portfolioserver.NewRouter(handler, a.Logger, reg)
	a.server = httpserver.NewServer(r, a.cfg.Server)

	grpcHandler := grpcv1.NewHandler(a.svc, a.Logger)
	grpcSrv, err := grpc.NewServer(a.cfg.Grpc.Addr(), grpcHandler, a.Logger)
	if err != nil {
		return fmt.Errorf("ошибка создания gRPC server: %w", err)
	}
	a.grpcSrv = grpcSrv

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- fmt.Errorf("http server остановлен с ошибкой: %w", err)
		}
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.grpcSrv.Start(); err != nil {
			a.errChan <- fmt.Errorf("gRPC server остановлен с ошибкой: %w", err)
		}
	}()

	if a.cfg.ServiceRegistry.ManagerBaseURL != "" {
		regHTTP, err := httpclient.OptionalMtlsHTTPClient(a.cfg.ServiceRegistry.Mtls)
		if err != nil {
			return err
		}
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			_ = registryinstance.Run(ctx, a.Logger, registryinstance.Config{
				ManagerBaseURL: a.cfg.ServiceRegistry.ManagerBaseURL,
				Secret:         a.cfg.ServiceRegistry.APISecret,
				ServiceName:    "quik-portfolio",
				InstanceID:     a.cfg.ServiceRegistry.InstanceID,
				GrpcPublicAddr: a.cfg.ServiceRegistry.GrpcPublicAddr,
				Interval:       time.Duration(a.cfg.ServiceRegistry.HeartbeatIntervalSec) * time.Second,
				HTTPClient:     regHTTP,
			})
		}()
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

	if a.server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = a.server.Shutdown(shutdownCtx)
	}

	if a.grpcSrv != nil {
		a.grpcSrv.Stop()
	}

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
