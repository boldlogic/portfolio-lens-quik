CREATE INDEX idx_quik_money_limits_load_date
ON quik.money_limits (load_date)

DROP INDEX IF EXISTS quik.money_limits.idx_quik_money_limits_load_date
