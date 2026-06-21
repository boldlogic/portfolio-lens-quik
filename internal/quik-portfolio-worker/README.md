# quik-portfolio-worker

Фоновый сервис roll-forward лимитов: копирование срезов на новую дату для `quik.money_limits`, `quik.security_limits`, `quik.security_limits_otc`.
## Запуск

```bash
go run ./cmd/quik-portfolio-worker
```

Конфиг по умолчанию: `configs/quik-portfolio-worker-config.yaml`.
