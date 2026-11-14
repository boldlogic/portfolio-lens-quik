package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/compose"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	app := compose.New()

	logger, err := app.Init()
	if err != nil {
		log.Fatalf("не удалось проинициализировать приложение: %v", err)
	}

	err = app.Start(ctx)
	if err != nil {
		logger.Fatal("не удалось запустить приложение", zap.Error(err))
	}

	app.Wait(ctx)
}
