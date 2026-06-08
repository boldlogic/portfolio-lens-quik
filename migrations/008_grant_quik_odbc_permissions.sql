-- +goose Up

GRANT SELECT, INSERT, UPDATE, DELETE ON quik.money_limits TO quik_odbc_writer;
GRANT SELECT, INSERT, UPDATE, DELETE ON quik.security_limits TO quik_odbc_writer;
GRANT SELECT, INSERT, UPDATE, DELETE ON quik.current_quotes TO quik_odbc_writer;
