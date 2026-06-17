package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/application"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	app, err := application.New()
	if err != nil {
		log.Fatalf("Не удалось создать приложение: %v", err)
	}

	if err = app.Start(ctx); err != nil {
		app.Logger.Fatal("Не удалось запустить приложение", zap.Error(err))
	}
	app.Logger.Info("приложение запущено")
	app.Wait(ctx)
	app.Logger.Info("приложение завершилось")
}
