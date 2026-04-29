package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/quik-currency/internal/application"
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
		log.Fatalf("Не удалось запустить приложение: %v", err)
	}
	err = app.Wait(ctx, cancel)
	if err != nil {
		log.Fatalf("Приложение завершилось с ошибкой: %v", err)
	}
}
