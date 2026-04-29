# portfolio-lens-quik

Репозиторий домена QUIK в рамках Portfolio Lens. Сейчас основной рабочий сервис - `quik-portfolio`; структура подготовлена для расширения на несколько микросервисов.

## Что внутри

- `quik-portfolio/` - сервис лимитов/портфеля (HTTP + gRPC).
- `quik-currency/` - задел под отдельный сервис домена валют/курсов.
- `pkg/` - общий код (транспорт, модели, интеграционные утилиты).
- `proto/` - protobuf-контракты и сгенерированный Go-код.
- `scripts/sql/` - bootstrap и DDL для MSSQL.
- `scripts/create-odbc-dsn.ps1` - создание 64-bit System DSN для ODBC-выгрузки.

## Быстрый старт (quik-portfolio)

1. Подготовить БД MSSQL:
  - выполнить `scripts/sql/bootstrap/create_database.sql`;
  - выполнить DDL из `scripts/sql/DDL/`;
  - при необходимости создать app-user через `scripts/sql/bootstrap/create_app_user.sql`.
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

Для экспорта из QUIK используйте БД. Текущее стандартное имя БД: `portfolio_lens_quik`.

Пример создания DSN:

```powershell
.\scripts\create-odbc-dsn.ps1 -Force -DbName portfolio_lens_quik -Dsn64 QuikPortfolioLocal_64 -PromptPassword
```

