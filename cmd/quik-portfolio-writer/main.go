package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-portfolio-writer/application"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	a, err := application.New()
	if err != nil {
		log.Fatalf("Не удалось создать приложение: %v", err)
	}

	err = a.Start(ctx)
	if err != nil {
		a.Logger.Fatal("не удалось запустить приложение", zap.Error(err))
	}
	a.Logger.Info("приложение запущено")
	err = a.Wait(ctx, cancel)
	if err != nil {
		a.Logger.Fatal("приложение завершилось с ошибкой", zap.Error(err))
	}
	a.Logger.Info("приложение завершилось")

}
