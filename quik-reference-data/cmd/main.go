package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/app"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	a := app.New()

	logger, err := a.Init()
	if err != nil {
		log.Fatalf("не удалось проинициализировать приложение, %v", err)
	}

	err = a.Start(ctx)
	if err != nil {
		logger.Fatal("не удалось запустить приложение", zap.Error(err))
	}

	err = a.InitFeatures(ctx)
	if err != nil {
		logger.Fatal("не удалось запустить приложение", zap.Error(err))
	}

	err = a.Wait(ctx, cancel)
	if err != nil {
		logger.Fatal("приложение завершилось с ошибкой", zap.Error(err))

	}

}
