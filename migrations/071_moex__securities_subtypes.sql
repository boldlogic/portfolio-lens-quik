-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.security_bonds (
    sec_id NVARCHAR(20) NOT NULL PRIMARY KEY,
    face_unit NVARCHAR(10) NULL,
    face_value DECIMAL(18,6) NULL,
    initial_face_value DECIMAL(18,6) NULL,
    coupon_value DECIMAL(18,6) NULL,
    coupon_percent DECIMAL(18,6) NULL,
    coupon_frequency INT NULL,
    coupon_date DATE NULL,
    maturity_date DATE NULL,
    days_to_redemption INT NULL,
    delivery_type NVARCHAR(50) NULL,
    settle_type NVARCHAR(50) NULL,
    calc_mode NVARCHAR(50) NULL,
    calc_mode_descr NVARCHAR(255) NULL,
    frequency NVARCHAR(50) NULL,
    program_registry_number NVARCHAR(100) NULL,
    CONSTRAINT FK_security_bonds_securities FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id) ON DELETE CASCADE
);

CREATE TABLE moex.security_derivatives (
    sec_id NVARCHAR(20) NOT NULL PRIMARY KEY,
    option_type NVARCHAR(50) NULL,
    strike DECIMAL(18,6) NULL,
    underlying_asset NVARCHAR(100) NULL,
    exec_type NVARCHAR(50) NULL,
    margin_style NVARCHAR(50) NULL,
    contract_name NVARCHAR(255) NULL,
    group_type NVARCHAR(50) NULL,
    initial_value DECIMAL(18,6) NULL,
    start_date_moex DATE NULL,
    expiration_clr_sess NVARCHAR(50) NULL,
    series_name NVARCHAR(100) NULL,
    CONSTRAINT FK_security_derivatives_securities FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id) ON DELETE CASCADE
);

CREATE TABLE moex.security_trading_profiles (
    sec_id NVARCHAR(20) NOT NULL PRIMARY KEY,
    trading_session NVARCHAR(50) NULL,
    morning_session BIT NULL,
    evening_session BIT NULL,
    time_course_int NVARCHAR(50) NULL,
    time_course_clr NVARCHAR(50) NULL,
    CONSTRAINT FK_security_trading_profiles_securities FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO moex.security_bonds (
    sec_id, face_unit, face_value, initial_face_value, coupon_value, coupon_percent,
    coupon_frequency, coupon_date, maturity_date, days_to_redemption,
    delivery_type, settle_type, calc_mode, calc_mode_descr, frequency, program_registry_number
)
SELECT
    sec_id, face_unit, face_value, initial_face_value, coupon_value, coupon_percent,
    coupon_frequency, coupon_date, maturity_date, days_to_redemption,
    delivery_type, settle_type, calc_mode, calc_mode_descr, frequency, program_registry_number
FROM moex.securities
WHERE group_code LIKE N'%bond%' OR type_code LIKE N'%bond%'
   OR maturity_date IS NOT NULL
   OR coupon_percent IS NOT NULL OR coupon_value IS NOT NULL OR coupon_frequency IS NOT NULL OR coupon_date IS NOT NULL
   OR days_to_redemption IS NOT NULL;

INSERT INTO moex.security_derivatives (
    sec_id, option_type, strike, underlying_asset, exec_type, margin_style,
    contract_name, group_type, initial_value, start_date_moex, expiration_clr_sess, series_name
)
SELECT
    sec_id, option_type, strike, underlying_asset, exec_type, margin_style,
    contract_name, group_type, initial_value, start_date_moex, expiration_clr_sess, series_name
FROM moex.securities
WHERE option_type IS NOT NULL OR strike IS NOT NULL OR underlying_asset IS NOT NULL
   OR type_code LIKE N'%option%' OR type_code LIKE N'%future%';

INSERT INTO moex.security_trading_profiles (
    sec_id, trading_session, morning_session, evening_session, time_course_int, time_course_clr
)
SELECT
    sec_id, trading_session, morning_session, evening_session, time_course_int, time_course_clr
FROM moex.securities
WHERE trading_session IS NOT NULL
   OR morning_session IS NOT NULL OR evening_session IS NOT NULL
   OR time_course_int IS NOT NULL OR time_course_clr IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
DROP VIEW IF EXISTS moex.instruments;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP COLUMN
    face_unit, face_value, initial_face_value, coupon_value, coupon_frequency, maturity_date;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP COLUMN
    coupon_percent, coupon_date, days_to_redemption, delivery_type, settle_type,
    calc_mode, calc_mode_descr, frequency, program_registry_number;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP COLUMN
    option_type, strike, underlying_asset, exec_type, margin_style, contract_name,
    group_type, initial_value, start_date_moex, expiration_clr_sess, series_name;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP COLUMN
    trading_session, time_course_int, time_course_clr, morning_session, evening_session;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC(N'
CREATE VIEW moex.instruments AS
SELECT
    s.sec_id,
    s.full_name,
    s.isin,
    s.short_name,
    s.reg_number,
    s.lat_name,
    s.currency_id,
    s.type_code,
    st.type_name,
    s.group_code,
    sg.group_name,
    st.iss_type_id,
    sg.iss_group_id,
    s.issue_date,
    s.lot_size,
    b.face_unit,
    b.initial_face_value,
    b.face_value,
    b.coupon_percent,
    b.coupon_value,
    b.coupon_frequency,
    b.maturity_date
FROM moex.securities s
LEFT JOIN moex.security_types st ON st.type_code = s.type_code
LEFT JOIN moex.security_groups sg ON sg.group_code = s.group_code
LEFT JOIN moex.security_bonds b ON b.sec_id = s.sec_id;
');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS moex.instruments;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities ADD
    face_unit NVARCHAR(10) NULL,
    face_value DECIMAL(18,6) NULL,
    initial_face_value DECIMAL(18,6) NULL,
    coupon_value DECIMAL(18,6) NULL,
    coupon_frequency INT NULL,
    maturity_date DATE NULL,
    coupon_percent DECIMAL(18,6) NULL,
    coupon_date DATE NULL,
    days_to_redemption INT NULL,
    delivery_type NVARCHAR(50) NULL,
    settle_type NVARCHAR(50) NULL,
    calc_mode NVARCHAR(50) NULL,
    calc_mode_descr NVARCHAR(255) NULL,
    frequency NVARCHAR(50) NULL,
    program_registry_number NVARCHAR(100) NULL,
    option_type NVARCHAR(50) NULL,
    strike DECIMAL(18,6) NULL,
    underlying_asset NVARCHAR(100) NULL,
    exec_type NVARCHAR(50) NULL,
    margin_style NVARCHAR(50) NULL,
    contract_name NVARCHAR(255) NULL,
    group_type NVARCHAR(50) NULL,
    initial_value DECIMAL(18,6) NULL,
    start_date_moex DATE NULL,
    expiration_clr_sess NVARCHAR(50) NULL,
    series_name NVARCHAR(100) NULL,
    trading_session NVARCHAR(50) NULL,
    time_course_int NVARCHAR(50) NULL,
    time_course_clr NVARCHAR(50) NULL,
    morning_session BIT NULL,
    evening_session BIT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE s SET
    s.face_unit = b.face_unit,
    s.face_value = b.face_value,
    s.initial_face_value = b.initial_face_value,
    s.coupon_value = b.coupon_value,
    s.coupon_percent = b.coupon_percent,
    s.coupon_frequency = b.coupon_frequency,
    s.coupon_date = b.coupon_date,
    s.maturity_date = b.maturity_date,
    s.days_to_redemption = b.days_to_redemption,
    s.delivery_type = b.delivery_type,
    s.settle_type = b.settle_type,
    s.calc_mode = b.calc_mode,
    s.calc_mode_descr = b.calc_mode_descr,
    s.frequency = b.frequency,
    s.program_registry_number = b.program_registry_number
FROM moex.securities s
INNER JOIN moex.security_bonds b ON b.sec_id = s.sec_id;

UPDATE s SET
    s.option_type = d.option_type,
    s.strike = d.strike,
    s.underlying_asset = d.underlying_asset,
    s.exec_type = d.exec_type,
    s.margin_style = d.margin_style,
    s.contract_name = d.contract_name,
    s.group_type = d.group_type,
    s.initial_value = d.initial_value,
    s.start_date_moex = d.start_date_moex,
    s.expiration_clr_sess = d.expiration_clr_sess,
    s.series_name = d.series_name
FROM moex.securities s
INNER JOIN moex.security_derivatives d ON d.sec_id = s.sec_id;

UPDATE s SET
    s.trading_session = t.trading_session,
    s.morning_session = t.morning_session,
    s.evening_session = t.evening_session,
    s.time_course_int = t.time_course_int,
    s.time_course_clr = t.time_course_clr
FROM moex.securities s
INNER JOIN moex.security_trading_profiles t ON t.sec_id = s.sec_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS moex.security_trading_profiles;
DROP TABLE IF EXISTS moex.security_derivatives;
DROP TABLE IF EXISTS moex.security_bonds;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC(N'
CREATE VIEW moex.instruments AS
SELECT
    s.sec_id,
    s.full_name,
    s.isin,
    s.short_name,
    s.reg_number,
    s.lat_name,
    s.currency_id,
    s.type_code,
    st.type_name,
    s.group_code,
    sg.group_name,
    st.iss_type_id,
    sg.iss_group_id,
    s.issue_date,
    s.lot_size,
    s.face_unit,
    s.initial_face_value,
    s.face_value,
    s.coupon_percent,
    s.coupon_value,
    s.coupon_frequency,
    s.maturity_date
FROM moex.securities s
LEFT JOIN moex.security_types st ON st.type_code = s.type_code
LEFT JOIN moex.security_groups sg ON sg.group_code = s.group_code;
');
-- +goose StatementEnd
