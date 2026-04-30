package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/infra"
	"go.uber.org/zap"
)

const (
	defaultConfigPath = "quik-reference-data/internal/configs/config.yaml"
	errChanBufSize    = 1
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	app, err := start(ctx)
	if err != nil {
		app.logger.Fatal("Не удалось запустить приложение", zap.Error(err))
	}
	err = app.wait(ctx, cancel)
	if err != nil {
		app.logger.Fatal("Приложение завершилось с ошибкой", zap.Error(err))

	}

}

type application struct {
	config  *infra.Config
	logger  *zap.Logger
	repo    *infra.Repository
	errChan chan error
	wg      sync.WaitGroup
}

func start(ctx context.Context) (*application, error) {

	a := application{}
	configPath := commonconfig.GetConfigPath(defaultConfigPath)

	//var err error

	conf, err := infra.LoadConfig(configPath)
	if err != nil {
		return &application{}, err
	}
	a.config = conf
	a.logger = logger.New(conf.Log)

	repo, err := infra.NewRepository(ctx, a.config.Db.GetDSN(), a.logger)
	if err != nil {
		return &application{}, err
	}
	a.repo = repo
	a.errChan = make(chan error, errChanBufSize)

	return &a, nil
}

func (a *application) wait(ctx context.Context, cancel context.CancelFunc) error {
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
