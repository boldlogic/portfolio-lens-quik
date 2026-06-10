# quik-currency

`quik-currency` отвечает за валютный справочник и кросс-курсы, которые приходят из QUIK через таблицы MSSQL.
## Запуск

Из корня репозитория:

```bash
go run ./quik-currency/cmd
```

По умолчанию конфиг читается из:

```text
quik-currency/internal/configs/config.yaml
```

Путь к конфигу можно передать тем же механизмом `commonconfig`, который используется в остальных сервисах репозитория.

## Конфигурация

В конфиге нужны блоки:

- `log` - уровень, формат и файл логов;
- `db` - подключение к MSSQL.

## Инициализация валют

Если `ref.currencies` пустая, сервис заполняет ее активными валютами из ISO 4217.

После этого сервис обновляет пустые `currency_name` по котировкам QUIK:

- источник: `quik.current_quotes`;
- фильтр: `class_code = 'CROSSRATE'`;
- коды `SUR`, `RUR`, `RUB` нормализуются в `RUB`;
- новые валюты из QUIK автоматически не добавляются.

Если валюта есть в `current_quotes`, но отсутствует в `ref.currencies`, она игнорируется на этапе заполнения имени.

## Перенос кросс-курсов

Фоновый воркер раз в `60*time.Second` переносит кросс-курсы из `quik.current_quotes` в `market.fx_cbr_rates`.

Основные правила:

- берутся строки `class_code = 'CROSSRATE'`;
- курс берется из `COALESCE(close_price, last_price)`;
- коды `RUR`, `SUR`, `RUB` не переносятся как базовая валюта;
- валюта должна существовать в `ref.currencies` (через `ref.external_codes` или `iso_char_code`);
- `quote_iso_code` для записей в `market.fx_cbr_rates` равен `643` (RUB);
- повторный запуск не создает дубли по ключу `(date, quote_iso_code, base_iso_code)`.

Если код из QUIK не найден в `ref.currencies`, курс не переносится, а в лог пишется `warn` с датой и кодом валюты.