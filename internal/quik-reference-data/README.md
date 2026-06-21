# quik-reference-data

`quik-reference-data` отвечает за ручные справочники QUIK. Сервис публикует HTTP API для фирм брокеров и периодически дополняет справочник фирм по данным лимитов из MSSQL.

## HTTP API

Фирмы:

- `GET /api/v1/quik/firms/` - список фирм;
- `POST /api/v1/quik/firms/` - создание фирмы;
- `GET /api/v1/quik/firms/{id}` - получение фирмы по id;
- `PATCH /api/v1/quik/firms/{id}` - изменение названия фирмы.

Ответ фирмы:

```json
{
  "id": 1,
  "firmCode": "MC0000000000",
  "firmName": "Брокер"
}
```

### Создание фирмы

`POST /api/v1/quik/firms/`

```json
{
  "firmCode": "MC0000000000",
  "firmName": "Брокер"
}
```

Правила входного JSON:

- `Content-Type` должен начинаться с `application/json`;
- тело запроса не больше 64 KiB;
- неизвестные поля не принимаются;
- `firmCode` обязателен, длина от 1 до 12 символов;
- `firmName` обязателен, длина от 1 до 128 символов.

При успешном создании сервис возвращает `201` и созданную фирму. Если фирма с таким `firmCode` уже есть, возвращается `409`.

### Изменение фирмы

`PATCH /api/v1/quik/firms/{id}`

```json
{
  "firmName": "Новое название"
}
```

Меняется только `firmName`. Поле `firmCode` через этот endpoint не изменяется. Если фирма не найдена, возвращается `404`.

## Ошибки HTTP

Ошибки возвращаются в JSON:

```json
{
  "title": "VALIDATION_ERROR",
  "status": 400,
  "detail": "некорректный id фирмы"
}
```

Используемые статусы:

- `400 VALIDATION_ERROR` - некорректный `id`, JSON или поля запроса;
- `404 NOT_FOUND` - фирма не найдена;
- `409 CONFLICT` - фирма с таким `firmCode` уже существует;
- `413 REQUEST_ENTITY_TOO_LARGE` - тело запроса превышает ограничение;
- `415 UNSUPPORTED_MEDIA_TYPE` - `Content-Type` не `application/json`;
- `500 SERVER_ERROR` - внутренняя ошибка, клиент получает `detail: "что-то пошло не так"`.

## Фоновая синхронизация фирм

Воркер `actualize_firms` раз в 60 секунд добавляет в `ref.firms` отсутствующие фирмы из лимитов:

- берет `firm_code` и очищенный от пробелов `firm_name`;
- читает данные из `quik.money_limits` и `quik.security_limits`;
- пропускает строки без `firm_code` и строки с пустым `firm_name`;
- не обновляет уже существующие фирмы, если `code` уже есть в `ref.firms`.

## Запуск

Из корня репозитория:

```bash
go run ./cmd/quik-reference-data
```

По умолчанию сервис ищет конфиг по пути:

```text
configs/quik-reference-data-config.yaml
```

## Конфигурация

В конфиге нужны блоки:

- `log` - настройки логирования;
- `server` - настройки HTTP server;
- `db` - подключение к MSSQL.
