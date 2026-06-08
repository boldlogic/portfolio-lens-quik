# portfolio-lens-quik

[CI](https://github.com/boldlogic/portfolio-lens-quik/actions/workflows/go.yml)

## Что внутри

- `quik-portfolio/` - сервис лимитов/портфеля (HTTP + gRPC).
- `quik-currency/` - сервис домена валют/курсов, см. `[quik-currency/README.md](quik-currency/README.md)`.
- `quik-reference-data/` - сервис справочников QUIK, см. `[quik-reference-data/README.md](quik-reference-data/README.md)`.
- `pkg/` - общий код (транспорт, модели, интеграционные утилиты).
- `proto/` - protobuf-контракты и сгенерированный Go-код.
- `scripts/sql/` - bootstrap и DDL для MSSQL.
- `scripts/create-odbc-dsn.ps1` - создание 64-bit System DSN для ODBC-выгрузки.

## Быстрый старт (quik-portfolio)

1. Подготовить БД MSSQL:
  - выполнить `scripts/sql/bootstrap/create_database.sql`;
  - выполнить DDL из `scripts/sql/DDL/`;
  - создать пользователей через `scripts/sql/bootstrap/create_*_user.sql` (reader, writer, worker, currency_worker) до миграции grants;
2. Проверить конфиг сервиса в `quik-portfolio/internal/configs/`.
3. Запустить сервис из корня репозитория:

```bash
go run ./quik-portfolio/cmd
go run ./quik-portfolio/cmd -config quik-portfolio/internal/configs/config.yaml
```

Порты по умолчанию:

- HTTP: `5030`
- gRPC: `5051`

## ODBC / QUIK

Для экспорта из QUIK: БД с миграциями, login `quik_odbc_writer` (только запись в `quik.*`, без `ref`/`market`), System DSN.

1. `scripts/sql/bootstrap/create_quik_odbc_user.sql` (до миграции `008`)
2. `go run ./migrator -command up` (если `008` ещё не применена)
3. DSN (cmd или PowerShell **от Administrator**; пароль и сервер из `.env`):

```cmd
.\scripts\create-odbc-dsn.cmd -Force
```

