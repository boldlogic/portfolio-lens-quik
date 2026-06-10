# quik-portfolio-worker

Фоновый сервис roll-forward лимитов: копирование срезов на новую дату для `quik.money_limits`, `quik.security_limits`, `quik.security_limits_otc`.
## Запуск

```bash
go run ./quik-portfolio-worker/cmd
```

Конфиг по умолчанию: `quik-portfolio-worker/internal/configs/config.yaml`.
