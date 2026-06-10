# quik-portfolio

**Назначение модуля:** read-only сервис: чтение лимитов из схемы `quik`, расчёт позиций портфеля с пересчётом в целевую валюту, выдача по HTTP и gRPC.

## HTTP API

| Метод | Путь | Назначение |
|-------|------|------------|
| GET | `/quik/money-limits` | денежные лимиты |
| GET | `/quik/security-limits` | бумажные лимиты (биржа) |
| GET | `/quik/otc-security-limits` | бумажные лимиты OTC |
| GET | `/quik/money-positions` | денежные позиции с оценкой |
| GET | `/quik/security-positions` | бумажные позиции с оценкой |
| GET | `/quik/otc-security-positions` | OTC-позиции с оценкой |
| GET | `/quik/positions` | сводные позиции (деньги + бумаги + OTC) |

## gRPC

Сервис `quikportfolio.v1.LimitsService`: read-only лимиты (денежные, бумажные, OTC). 

## Запуск

Предварительно: развёрнута БД с миграциями из `migrations/`, заполнен конфиг. 
Из корня репозитория:

```bash
go run ./quik-portfolio/cmd
go run ./quik-portfolio/cmd -config path/to/config.yaml
```

Путь к YAML по умолчанию: `quik-portfolio/internal/configs/config.yaml`.
