# quik-portfolio

**Назначение модуля:** приём и хранение лимитов из контура QUIK (схема `quik` в SQL Server), расчёт портфеля с пересчётом оценки в целевую валюту, выдача данных по HTTP и gRPC; справочник фирм.

**Интерфейсы:** спецификация OpenAPI ведётся как внешний артефакт; проектное описание [docs/architecture.md](../docs/architecture.md).

**Сеть по умолчанию (Docker-пример):** HTTP **5030** (`server.port`), gRPC **5051** (`grpc.port`).

## Запуск

Предварительно: развёрнута БД с DDL из `scripts/sql/DDL/`, заполнен конфиг. Секреты БД задаются через переменные окружения (секция `db` в конфиге); в Docker часто достаточно `MSSQL_SA_PASSWORD` в `.env`.

Из корня репозитория:

```bash
go run ./quik-portfolio/cmd
go run ./quik-portfolio/cmd -config path/to/config.yaml
```

Путь к YAML по умолчанию: `quik-portfolio/internal/configs/config.yaml`.
