package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/application"
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
	log.Println("приложение запущено")

	err = app.Wait(ctx, cancel)
	if err != nil {
		log.Fatalf("приложение завершилось с ошибкой: %v", err)
	}
	log.Println("приложение завершилось без ошибок")
}
