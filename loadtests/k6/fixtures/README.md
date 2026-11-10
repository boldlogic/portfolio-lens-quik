# Фикстуры для k6 prepare

## firms.json

Список фирм для prepare лимитов. **Новые фирмы через API не создаются** (в `ref.firms` только `id tinyint`, места мало).

Поля:

| Поле | Обязательно | Описание |
|------|-------------|----------|
| `firmCode` | да | до 12 символов, как в `ref.firms.code` |
| `firmName` | нет | только для справки в файле |
| `id` | нет | для read лимитов не нужен. Для mixed/get_firm — если есть, prepare не дергает GET списка |

Коды должны уже быть в БД (QUIK, ручной ввод, sync worker). Перед прогоном сверьте:

```sql
SELECT id, code, name FROM ref.firms WHERE code IN (...);
```
