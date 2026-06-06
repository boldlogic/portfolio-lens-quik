# quik-portfolio

**Назначение модуля:** read-only сервис: чтение лимитов из схемы `quik`, расчёт портфеля с пересчётом в целевую валюту, выдача по HTTP и gRPC.

## Запуск

Предварительно: развёрнута БД с DDL из `scripts/sql/DDL/`, заполнен конфиг. Секреты БД задаются через переменные окружения (секция `db` в конфиге).

Из корня репозитория:

```bash
go run ./quik-portfolio/cmd
go run ./quik-portfolio/cmd -config path/to/config.yaml
```

Путь к YAML по умолчанию: `quik-portfolio/internal/configs/config.yaml`.
