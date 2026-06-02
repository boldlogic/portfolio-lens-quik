package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/application"
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
		log.Fatalf("не удалось создать приложение: %v", err)
	}

	if err = app.Start(ctx); err != nil {
		log.Fatalf("не удалось запустить приложение: %v", err)
	}

	app.Logger.Info("приложение запущено")
	err = app.Wait()
	if err != nil {
		app.Logger.Fatal("приложение завершилось с ошибкой", zap.Error(err))
	}
	app.Logger.Info("приложение завершилось без ошибок")
}
