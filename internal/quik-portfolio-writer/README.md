# quik-portfolio-writer

HTTP-сервис записи лимитов: денежные, биржевые бумаги, OTC. Fallback и ручная загрузка данных, которые QUIK не выгружает по ODBC (OTC) или при отсутствии у брокера QUIK.

## HTTP API

Префикс `/api/v1`.

| Метод | Путь |
|-------|------|
| POST | `/money-limits` |
| POST | `/security-limits` |
| POST | `/otc-security-limits` |

## Таблицы

Запись в `quik.money_limits`, `quik.security_limits`, `quik.security_limits_otc`. Справочник фирм: join `ref.firms` по `firm_code`.

## Запуск

```bash
go run ./quik-portfolio-writer/cmd
```

Конфиг по умолчанию: `quik-portfolio-writer/internal/configs/config.yaml`.
