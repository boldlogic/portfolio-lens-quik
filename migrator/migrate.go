package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbconfig"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/pressly/goose/v3"
)

const (
	defaultConfigPath = "migrator/config.yaml"
	defaultCommand    = "up"
	migrationsDir     = "migrations"
	gooseDialect      = "mssql"
)

type dbConfig = dbconfig.DBConfig

func main() {
	command := flag.String("command", defaultCommand, "команда goose (пока только up)")
	configPath := flag.String("config", defaultConfigPath, "путь к YAML с параметрами БД")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := openDB(cfg.GetDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.SetDialect(gooseDialect); err != nil {
		log.Fatal(err)
	}

	if err := runCommand(db, *command); err != nil {
		log.Fatal(err)
	}
}

func loadConfig(path string) (*dbConfig, error) {
	cfg, err := commonconfig.DecodeConfigStrict[dbConfig](path)
	if err != nil {
		return nil, err
	}
	cfg.ApplyDefaults()
	if errs := cfg.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("некорректный конфиг: %w", errors.Join(errs...))
	}
	return &cfg, nil
}

func openDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("некорректный DSN")
	}
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть подключение к БД: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("не удалось проверить подключение к БД: %w", err)
	}
	return db, nil
}

func runCommand(db *sql.DB, command string) error {
	switch command {
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			return fmt.Errorf("ошибка при миграции: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("неизвестная команда: %q", command)
	}
}
