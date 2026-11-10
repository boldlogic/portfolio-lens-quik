# k6 нагрузочные тесты (Portfolio Lens QUIK)

Сверка путей API: 2026-06-10. Reader: `/api/v1/quik/*`. Writer: `/api/v1/money-limits` (без `/quik`). Firms: reference-data `:6030`.

Нагрузка по HTTP для трёх сервисов:

| Сервис | Дефолт URL | Роль в профиле |
|--------|------------|----------------|
| quik-portfolio | `http://127.0.0.1:5030` | чтение лимитов и портфеля |
| quik-portfolio-writer | `http://127.0.0.1:5060` | запись лимитов |
| quik-reference-data | `http://127.0.0.1:6030` | справочник фирм |

Порты взяты из локальных `config.yaml` (`server.listen_port`). Если writer слушает другой порт (например 5060), передайте `-WriterBaseUrl`.

## Предусловия

1. MSSQL `portfolio_lens_quik` доступна (как в конфигах сервисов).
2. Запущены три процесса, например из корня репозитория:

```bash
go run ./quik-reference-data/cmd
go run ./quik-portfolio-writer/cmd
go run ./quik-portfolio/cmd
```

3. Установлен [k6](https://grafana.com/docs/k6/latest/set-up/install-k6/).

## Быстрый запуск

```powershell
.\loadtests\k6\run.ps1
```

То же с явным сценарием:

```powershell
.\loadtests\k6\run.ps1 -Scenario mixed
```

При блокировке политики PowerShell:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\loadtests\k6\run.ps1
```

## Дефолты

- `TargetRps = 1000`
- длительность по умолчанию: `DurationSec = 30` в каждом `run_*.ps1` (разогрев до `5s` + удержание; минимум `6`)
- prepare: фирмы из `loadtests/k6/fixtures/firms.json`, `20` client codes; тип лимита зависит от сценария (см. ниже)
- query `loadDate` по умолчанию **не передаётся** (сервер подставляет «сегодня»); явная дата: `-IncludeDate`
- отчёты: `docs/reports/k6`
- state: `loadtests/k6/state` (генерируется prepare, в git не коммитится)

## Итоговый mixed-профиль

GET лимитов в mixed **нет** (отдельные сценарии).

| Endpoint | Доля RPS |
|----------|----------:|
| `GET /api/v1/quik/positions?targetCurrency=RUB` | 55% |
| `GET /api/v1/quik/firms/` | 12% |
| `GET /api/v1/quik/firms/{id}` | 5% |
| `PUT /api/v1/limits` (money) | 10% |
| `PUT /api/v1/limits` (security) | 10% |
| `PUT /api/v1/limits` (security_otc) | 4% |
| `POST /api/v1/quik/firms/` | 2% |
| `PATCH /api/v1/quik/firms/{id}` | остаток |

`GET /health` и `GET /metrics` в нагрузку не входят.

Примеры:

```powershell
.\loadtests\k6\run.ps1 -TargetRps 100
.\loadtests\k6\run_mixed.ps1 -TargetRps 500
.\loadtests\k6\run_get_firms.ps1 -TargetRps 10
```

## Лимиты: by_client и paged

Для `money`, `securities`, `securities/otc` два режима:

| Режим | Query | k6-итерация |
|-------|-------|-------------|
| **by_client** | один `clientCodes`, без `limit`/`offset`/`loadDate`/`includeTotalCount` (дефолт) | один GET |
| **paged** | `limit` + `offset` + `includeTotalCount=true`, без `clientCodes`/`loadDate` (дефолт) | полный обход страниц (несколько GET) |

Ответ read-лимитов: JSON-объект `{ limits, limit, offset, totalCount? }` (не массив).

`run_get_limits.ps1`: за одну итерацию k6 три GET (`money-limits`, `security-limits`, `otc-security-limits`) с одним `clientCodes`; prepare пишет все три типа лимитов через writer.

### quik-portfolio (read лимиты)

| Endpoint | by_client | paged |
|----------|-----------|-------|
| `GET .../money-limits` | `run_get_money_limits_by_client.ps1` | `run_get_money_limits_paged.ps1` |
| `GET .../security-limits` | `run_get_security_limits_by_client.ps1` | `run_get_security_limits_paged.ps1` |
| `GET .../otc-security-limits` | `run_get_security_limits_otc_by_client.ps1` | `run_get_security_limits_otc_paged.ps1` |
| три read-эндпоинта | `run_get_limits.ps1` (3 GET/iter) | — |

```powershell
# Limit read: TargetRps v run_get_*.ps1 = HTTP GET/s
# by_client: 1 GET/iter. get_limits: 3 GET/iter. paged: iter/s = HTTP / ceil(totalCount/page_limit)
# DurationSec: полная длительность теста в секундах (в каждом run_*.ps1 свой дефолт)
.\loadtests\k6\run_get_money_limits_paged.ps1
.\loadtests\k6\run_get_money_limits_by_client.ps1

.\loadtests\k6\run_get_money_limits_by_client.ps1 -TargetRps 10 -DurationSec 60
.\loadtests\k6\run_get_limits.ps1 -TargetRps 100 -DurationSec 45
```

Эквивалент через `run.ps1`:

```powershell
.\loadtests\k6\run.ps1 -Scenario get_money_limits_by_client -TargetRps 10
```

### quik-portfolio (positions)

Ответ: объект `{ marketValueTotal, targetCurrency, positions[] }` (не массив). Query: `targetCurrency` (дефолт RUB на сервере), опционально `clientCodes`.

| Endpoint | Сценарий | Запуск |
|----------|----------|--------|
| `GET .../positions` | `get_positions` | `.\loadtests\k6\run_get_positions.ps1` |
| `GET .../money-positions` | `get_money_positions` | `.\loadtests\k6\run_get_money_positions.ps1` |
| `GET .../security-positions` | `get_security_positions` | `.\loadtests\k6\run_get_security_positions.ps1` |
| `GET .../otc-security-positions` | `get_otc_security_positions` | `.\loadtests\k6\run_get_otc_security_positions.ps1` |

`run_get_portfolio.ps1` — обёртка на `run_get_positions.ps1` (устаревший путь `/portfolios`).

```powershell
.\loadtests\k6\run_get_positions.ps1 -TargetRps 10 -DurationSec 30
.\loadtests\k6\run_get_security_positions.ps1 -TargetRps 10 -DurationSec 6
```

### quik-portfolio-writer (write)

Writer: `PUT /api/v1/limits`, `POST /api/v1/limits/upload` (CSV `;`, multipart). Заголовок `X-API-KEY` обязателен. Успех: **204 No Content**.

| Endpoint | Сценарий | Запуск |
|----------|----------|--------|
| `PUT /api/v1/limits` (money) | `post_money_limit` | `.\loadtests\k6\run_post_money_limit.ps1` |
| `PUT /api/v1/limits` (security) | `post_security_limit` | `.\loadtests\k6\run_post_security_limit.ps1` |
| `PUT /api/v1/limits` (security_otc) | `post_security_limit_otc` | `.\loadtests\k6\run_post_security_limit_otc.ps1` |
| 50% PUT + 50% CSV (1-2 строки) | `write_limits_mixed` | `.\loadtests\k6\run_write_limits_mixed.ps1` |

```powershell
.\loadtests\k6\run_write_limits_mixed.ps1 `
  -WriterBaseUrl http://127.0.0.1:5030 `
  -WriterApiKey ksjdfgld `
  -TargetRps 50 `
  -DurationSec 30
```

`-WriterApiKey` или env `WRITER_API_KEY` для writer-сценариев.

### quik-reference-data

| Endpoint | Сценарий | Запуск |
|----------|----------|--------|
| `GET /api/v1/quik/firms/` | `get_firms` | `.\loadtests\k6\run_get_firms.ps1` |
| `GET /api/v1/quik/firms/{id}` | `get_firm` | `.\loadtests\k6\run_get_firm.ps1` |
| `POST /api/v1/quik/firms/` | `post_firm` | `.\loadtests\k6\run_post_firm.ps1` |
| `PATCH /api/v1/quik/firms/{id}` | `patch_firm` | `.\loadtests\k6\run_patch_firm.ps1` |

Смешанная нагрузка только в `mixed` / `run_mixed.ps1`.

Параметры URL и prepare:

```powershell
.\loadtests\k6\run.ps1 `
  -PortfolioBaseUrl http://127.0.0.1:5030 `
  -WriterBaseUrl http://127.0.0.1:5060 `
  -ReferenceDataBaseUrl http://127.0.0.1:6030 `
  -PrepareClientCount 20 `
  -SkipPrepare
```

`-SkipPrepare` пропускает HTTP prepare; `target_rps` в state перезаписывается из `-TargetRps` (иначе в JSON мог остаться старый 1000).

Цель нагрузки: для **`run_get_*` лимитов** и **`run_get_limits.ps1`** параметр **`TargetRps` = HTTP GET/s** (в каждом `run_get_*.ps1`). Перед k6: probe portfolio; если не отвечает, PowerShell завершает запуск с ошибкой. Во время теста при connectex: `test.abort` + порог `http_req_failed` (`rate<0.01`, `delayAbortEval: 0s`). `test.abort` не отменяет уже запущенные `http.get`, поэтому timeout GET 3s и `gracefulStop: 0s` (в логе k6 должно быть `gracefulStop: 0s`, не 30s). Для `*_paged` k6 iter/s = HTTP / число страниц. Остальные сценарии (`mixed`, POST, firms): `-TargetRps` в `run.ps1` = iter/s как раньше.

Обёртки `run_*.ps1` передают `TargetRps` в дочерний скрипт через **явный hashtable** (`TargetRps = $TargetRps`), не через `@PSBoundParameters`: иначе дефолт из `param()` родителя не попадает в `run_limit_test.ps1` / `run.ps1` и в отчёт уходит `rps1`. В консоли cyan-строка с фактическим `TargetRps` печатается из `run.ps1`. Запускайте `.ps1`, не `scenarios/*.js` (см. `lib/launchGuard.js`).

## Подготовка данных

Перед k6 скрипт `prepare/<scenario>.ps1`:

1. Фирмы из [`fixtures/firms.json`](fixtures/firms.json): **БКС** (`NC0058900000`, id=1) и **ВТБ** (`MC0003300000`, id=2). Синтетические `LD...` удалены.
2. Коды лимитов из [`fixtures/limit_codes.json`](fixtures/limit_codes.json): writer POST с полями `currency` / `secCode` (в read positions то же значение в поле `ticker`).
3. `clientCode` — уникальные `LDCL...` на каждый prepare (`Format-LoadCode`), не коды из prod-примера.
4. Сценарии **read лимитов / positions** (`get_*_limits_*`, `get_*_positions`): `-FixtureFirmsCodesOnly`, reference-data не нужен. Для `mixed`, `get_firm` — `id` из `firms.json`.
5. Writer seed по `-WriterSeedKinds`:
   - money-only: `POST /api/v1/money-limits` (`currency` из фикстуры);
   - securities / OTC: свой endpoint, `secCode` из фикстуры;
   - `get_limits` / `get_positions` / `mixed`: `All`.
6. **Котировки `market.*`** k6 не сидит — заполнение вручную в БД. Пустой `positions: []` при 200 допустим для smoke.
7. `primary_client_code` = первый из `client_codes` (для by_client limits).
8. Запись `loadtests/k6/state/<scenario>.json`.

```powershell
.\loadtests\k6\run_get_money_limits_by_client.ps1 -PrepareClientCount 20 -WriterBaseUrl http://127.0.0.1:5060
# для read лимитов достаточно writer + portfolio; reference-data на 6030 не нужен
.\loadtests\k6\run.ps1 -Scenario get_security_limits_paged -FirmsFixturePath .\my-firms.json
```

`-PrepareFirmCount 0` — все фирмы из файла. `N` — первые N записей.

Проверка money-only prepare: в логе writer только POST `/api/v1/money-limits`, без securities/otc.

## Ручная проверка одного запроса

```powershell
$ref = "http://127.0.0.1:6030"
Invoke-RestMethod -Method Get -Uri "$ref/api/v1/quik/firms/" -Headers @{ Accept = "application/json" }

$writer = "http://127.0.0.1:5030"
$body = @{
  limitType = "money"
  clientCode = "TESTCLIENT01"
  ticker = "RUB"
  positionCode = "EQTV"
  settleCode = "Tx"
  firmCode = "MC0000000000"
  balance = "1.00"
} | ConvertTo-Json
Invoke-RestMethod -Method Put -Uri "$writer/api/v1/limits" -ContentType "application/json" -Headers @{ "X-API-KEY" = "your-api-key" } -Body $body
```

`firmCode` должен существовать в справочнике фирм.

## Результаты

```text
docs/reports/k6/load_<scenario>__dur<DurationSec>s__rps<target>__<pass|fail>__<local-time>.html
docs/reports/k6/last_<scenario>.html
```
