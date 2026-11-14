-- +goose Up
IF EXISTS (
    SELECT 1 FROM sys.foreign_keys
    WHERE referenced_object_id = OBJECT_ID(N'ref.boards')
      AND name NOT IN (N'FK_ref_instrument_type_boards_board_id', N'FK_ref_instrument_boards_board', N'FK_moex_common_board_map_board', N'FK_market_moex_quotes_board', N'FK_market_moex_history_board')
)
    THROW 51030, N'Найдена неучтённая внешняя ссылка на ref.boards; расширение board_id остановлено', 1;

ALTER TABLE ref.instrument_type_boards DROP CONSTRAINT FK_ref_instrument_type_boards_board_id;
ALTER TABLE ref.instrument_type_boards DROP CONSTRAINT UQ_ref_instrument_type_boards_board_id;
ALTER TABLE ref.instrument_type_boards DROP CONSTRAINT PK_ref_instrument_type_boards;
ALTER TABLE ref.instrument_boards DROP CONSTRAINT FK_ref_instrument_boards_board;
ALTER TABLE ref.instrument_boards DROP CONSTRAINT PK_ref_instrument_boards;
ALTER TABLE moex.common_board_map DROP CONSTRAINT FK_moex_common_board_map_board;
ALTER TABLE market.moex_quotes DROP CONSTRAINT FK_market_moex_quotes_board;
ALTER TABLE market.moex_quotes DROP CONSTRAINT PK_market_moex_quotes;
ALTER TABLE market.moex_historical_rates DROP CONSTRAINT FK_market_moex_history_board;
ALTER TABLE market.moex_historical_rates DROP CONSTRAINT PK_market_moex_historical_rates;
ALTER TABLE ref.boards DROP CONSTRAINT FK_ref_boards_trade_point_id;
ALTER TABLE ref.boards DROP CONSTRAINT PK_ref_boards;
DROP INDEX UQ_boards_code ON ref.boards;

ALTER TABLE ref.instrument_type_boards ALTER COLUMN board_id SMALLINT NOT NULL;
ALTER TABLE ref.instrument_boards ALTER COLUMN board_id SMALLINT NOT NULL;
ALTER TABLE moex.common_board_map ALTER COLUMN board_id SMALLINT NOT NULL;
ALTER TABLE market.moex_quotes ALTER COLUMN board_id SMALLINT NOT NULL;
ALTER TABLE market.moex_historical_rates ALTER COLUMN board_id SMALLINT NOT NULL;

CREATE TABLE ref.boards_next (
    board_id SMALLINT IDENTITY(1,1) NOT NULL,
    code NVARCHAR(12) NOT NULL,
    name NVARCHAR(100) NOT NULL,
    trade_point_id TINYINT NULL,
    is_traded BIT NOT NULL CONSTRAINT DF_ref_boards_next_is_traded DEFAULT (0),
    CONSTRAINT PK_ref_boards_next PRIMARY KEY CLUSTERED (board_id),
    CONSTRAINT FK_ref_boards_next_trade_point FOREIGN KEY (trade_point_id) REFERENCES ref.trade_points(point_id)
);
SET IDENTITY_INSERT ref.boards_next ON;
INSERT INTO ref.boards_next (board_id, code, name, trade_point_id, is_traded)
SELECT board_id, code, name, trade_point_id, is_traded FROM ref.boards;
SET IDENTITY_INSERT ref.boards_next OFF;
DROP TABLE ref.boards;
EXEC sys.sp_rename N'ref.boards_next', N'boards';
CREATE UNIQUE NONCLUSTERED INDEX UQ_boards_code ON ref.boards(code);

ALTER TABLE ref.instrument_type_boards ADD CONSTRAINT PK_ref_instrument_type_boards PRIMARY KEY CLUSTERED (type_id, board_id);
ALTER TABLE ref.instrument_type_boards ADD CONSTRAINT UQ_ref_instrument_type_boards_board_id UNIQUE (board_id);
ALTER TABLE ref.instrument_type_boards ADD CONSTRAINT FK_ref_instrument_type_boards_board_id FOREIGN KEY (board_id) REFERENCES ref.boards(board_id);
ALTER TABLE ref.instrument_boards ADD CONSTRAINT PK_ref_instrument_boards PRIMARY KEY CLUSTERED (instrument_id, board_id);
ALTER TABLE ref.instrument_boards ADD CONSTRAINT FK_ref_instrument_boards_board FOREIGN KEY (board_id) REFERENCES ref.boards(board_id);
ALTER TABLE moex.common_board_map ADD CONSTRAINT FK_moex_common_board_map_board FOREIGN KEY (board_id) REFERENCES ref.boards(board_id);
ALTER TABLE market.moex_quotes ADD CONSTRAINT PK_market_moex_quotes PRIMARY KEY CLUSTERED (instrument_id, board_id, quote_date);
ALTER TABLE market.moex_quotes ADD CONSTRAINT FK_market_moex_quotes_board FOREIGN KEY (board_id) REFERENCES ref.boards(board_id);
ALTER TABLE market.moex_historical_rates ADD CONSTRAINT PK_market_moex_historical_rates PRIMARY KEY CLUSTERED (instrument_id, board_id, trade_date);
ALTER TABLE market.moex_historical_rates ADD CONSTRAINT FK_market_moex_history_board FOREIGN KEY (board_id) REFERENCES ref.boards(board_id);

IF NOT EXISTS (SELECT 1 FROM ref.instrument_types WHERE title = N'Валютные пары')
    INSERT INTO ref.instrument_types (title) VALUES (N'Валютные пары');
IF NOT EXISTS (SELECT 1 FROM ref.instrument_types WHERE title = N'Фьючерсы')
    INSERT INTO ref.instrument_types (title) VALUES (N'Фьючерсы');
IF NOT EXISTS (SELECT 1 FROM ref.instrument_types WHERE title = N'Опционы')
    INSERT INTO ref.instrument_types (title) VALUES (N'Опционы');
IF NOT EXISTS (SELECT 1 FROM ref.instrument_types WHERE title = N'Индексы')
    INSERT INTO ref.instrument_types (title) VALUES (N'Индексы');

CREATE TABLE ref.instrument_subtypes (
    subtype_id SMALLINT IDENTITY(1,1) NOT NULL,
    type_id TINYINT NOT NULL,
    code VARCHAR(40) NOT NULL,
    title NVARCHAR(100) NOT NULL,
    CONSTRAINT PK_ref_instrument_subtypes PRIMARY KEY (subtype_id),
    CONSTRAINT FK_ref_instrument_subtypes_type FOREIGN KEY (type_id) REFERENCES ref.instrument_types(type_id),
    CONSTRAINT UQ_ref_instrument_subtypes_type_code UNIQUE (type_id, code),
    CONSTRAINT UQ_ref_instrument_subtypes_type_subtype UNIQUE (type_id, subtype_id)
);

-- +goose StatementBegin
DECLARE @shares_type_id TINYINT = (SELECT type_id FROM ref.instrument_types WHERE title = N'Акции');
IF @shares_type_id IS NOT NULL
    INSERT INTO ref.instrument_subtypes (type_id, code, title)
    VALUES (@shares_type_id, 'ETF', N'ETF'), (@shares_type_id, 'MUTUAL_FUND', N'ПИФ');
-- +goose StatementEnd

ALTER TABLE ref.instruments ADD subtype_id SMALLINT NULL;
ALTER TABLE ref.instruments ADD description NVARCHAR(1000) NULL;
ALTER TABLE ref.instruments ALTER COLUMN registration_number NVARCHAR(64) NULL;
ALTER TABLE ref.instruments ALTER COLUMN full_name NVARCHAR(255) NULL;
ALTER TABLE ref.instruments ALTER COLUMN short_name NVARCHAR(200) NULL;
ALTER TABLE ref.instruments ADD CONSTRAINT FK_ref_instruments_subtype FOREIGN KEY (type_id, subtype_id)
    REFERENCES ref.instrument_subtypes(type_id, subtype_id);

CREATE TABLE ref.instrument_bond_details (
    instrument_id BIGINT NOT NULL,
    face_value DECIMAL(19,8) NULL,
    initial_face_value DECIMAL(19,8) NULL,
    current_face_value DECIMAL(19,8) NULL,
    maturity_date DATE NULL,
    coupon_duration INT NULL,
    coupon_rate DECIMAL(19,8) NULL,
    coupons_per_year DECIMAL(19,8) NULL,
    coupon_type NVARCHAR(100) NULL,
    next_coupon_date DATE NULL,
    issue_date DATE NULL,
    nominal_currency NVARCHAR(16) NULL,
    indexed_nominal BIT NULL,
    amortization BIT NULL,
    floating_coupon BIT NULL,
    delivery_type NVARCHAR(50) NULL,
    settlement_type NVARCHAR(50) NULL,
    calculation_mode NVARCHAR(50) NULL,
    calculation_description NVARCHAR(255) NULL,
    CONSTRAINT PK_ref_instrument_bond_details PRIMARY KEY (instrument_id),
    CONSTRAINT FK_ref_instrument_bond_details_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments(instrument_id)
);

CREATE TABLE ref.instrument_fund_details (
    instrument_id BIGINT NOT NULL,
    benchmark NVARCHAR(255) NULL,
    description NVARCHAR(1000) NULL,
    benchmark_description NVARCHAR(1000) NULL,
    total_expense DECIMAL(19,8) NULL,
    hurdle_rate DECIMAL(19,8) NULL,
    performance_fee DECIMAL(19,8) NULL,
    fixed_commission DECIMAL(19,8) NULL,
    expense_commission DECIMAL(19,8) NULL,
    buy_premium DECIMAL(19,8) NULL,
    sell_discount DECIMAL(19,8) NULL,
    watermark BIT NULL,
    CONSTRAINT PK_ref_instrument_fund_details PRIMARY KEY (instrument_id),
    CONSTRAINT FK_ref_instrument_fund_details_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments(instrument_id)
);

CREATE TABLE ref.instrument_future_details (
    instrument_id BIGINT NOT NULL,
    underlying_instrument_id BIGINT NULL,
    underlying_code NVARCHAR(100) NULL,
    contract_name NVARCHAR(255) NULL,
    expiration_date DATE NULL,
    start_date DATE NULL,
    lot_size DECIMAL(19,8) NULL,
    initial_value DECIMAL(19,8) NULL,
    CONSTRAINT PK_ref_instrument_future_details PRIMARY KEY (instrument_id),
    CONSTRAINT FK_ref_instrument_future_details_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments(instrument_id),
    CONSTRAINT FK_ref_instrument_future_details_underlying FOREIGN KEY (underlying_instrument_id) REFERENCES ref.instruments(instrument_id)
);

CREATE TABLE ref.instrument_option_details (
    instrument_id BIGINT NOT NULL,
    underlying_instrument_id BIGINT NULL,
    underlying_code NVARCHAR(100) NULL,
    option_type NVARCHAR(50) NULL,
    strike DECIMAL(19,8) NULL,
    expiration_date DATE NULL,
    exercise_type NVARCHAR(50) NULL,
    margin_style NVARCHAR(50) NULL,
    contract_name NVARCHAR(255) NULL,
    series_name NVARCHAR(255) NULL,
    calculation_mode NVARCHAR(50) NULL,
    calculation_description NVARCHAR(255) NULL,
    settlement_type NVARCHAR(50) NULL,
    CONSTRAINT PK_ref_instrument_option_details PRIMARY KEY (instrument_id),
    CONSTRAINT FK_ref_instrument_option_details_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments(instrument_id),
    CONSTRAINT FK_ref_instrument_option_details_underlying FOREIGN KEY (underlying_instrument_id) REFERENCES ref.instruments(instrument_id)
);

INSERT INTO ref.instrument_bond_details (instrument_id, face_value, maturity_date, coupon_duration)
SELECT instrument_id, face_value, maturity_date, coupon_duration FROM ref.instruments
WHERE type_id = (SELECT type_id FROM ref.instrument_types WHERE title = N'Облигации')
  AND (face_value IS NOT NULL OR maturity_date IS NOT NULL OR coupon_duration IS NOT NULL);

CREATE TABLE ref.instrument_attribute_priorities (
    attribute_code VARCHAR(80) NOT NULL,
    ext_system_id TINYINT NOT NULL,
    priority SMALLINT NOT NULL,
    CONSTRAINT PK_ref_instrument_attribute_priorities PRIMARY KEY (attribute_code, ext_system_id),
    CONSTRAINT FK_ref_instrument_attribute_priorities_system FOREIGN KEY (ext_system_id) REFERENCES ref.external_systems(ext_system_id)
);

INSERT INTO ref.instrument_attribute_priorities (attribute_code, ext_system_id, priority)
SELECT p.attribute_code, es.ext_system_id, p.priority
FROM (VALUES
    ('classification.new_card', 'MOEX', 100), ('classification.new_card', 'QUIK', 90), ('classification.new_card', 'TINVEST', 80), ('classification.new_card', 'BCS', 70),
    ('classification.subtype', 'MOEX', 100), ('classification.subtype', 'QUIK', 90), ('classification.subtype', 'TINVEST', 80), ('classification.subtype', 'BCS', 70),
    ('external_code', 'MOEX', 100), ('external_code', 'QUIK', 90), ('external_code', 'TINVEST', 80), ('external_code', 'BCS', 70),
    ('board.name', 'MOEX', 100), ('board.name', 'QUIK', 90),
    ('instrument_board.currency_id', 'MOEX', 100), ('instrument_board.currency_id', 'BCS', 90), ('instrument_board.currency_id', 'QUIK', 80),
    ('instrument_board.base_currency_id', 'QUIK', 100), ('instrument_board.quote_currency_id', 'QUIK', 100), ('instrument_board.counter_currency_id', 'QUIK', 100),
    ('instrument_board.is_primary', 'MOEX', 100), ('instrument_board.is_primary', 'BCS', 90),
    ('bond.face_value', 'MOEX', 100), ('bond.face_value', 'QUIK', 90), ('bond.face_value', 'BCS', 80), ('bond.face_value', 'TINVEST', 70),
    ('bond.initial_face_value', 'MOEX', 100),
    ('bond.current_face_value', 'TINVEST', 100),
    ('bond.nominal_currency', 'MOEX', 100), ('bond.nominal_currency', 'BCS', 90), ('bond.nominal_currency', 'TINVEST', 80),
    ('bond.indexed_nominal', 'TINVEST', 100), ('bond.amortization', 'TINVEST', 100), ('bond.floating_coupon', 'TINVEST', 100),
    ('bond.delivery_type', 'MOEX', 100), ('bond.settlement_type', 'MOEX', 100), ('bond.calculation_mode', 'MOEX', 100), ('bond.calculation_description', 'MOEX', 100),
    ('bond.coupon_type', 'BCS', 100), ('bond.issue_date', 'MOEX', 100), ('bond.issue_date', 'BCS', 90),
    ('bond.maturity_date', 'MOEX', 100), ('bond.maturity_date', 'BCS', 90), ('bond.maturity_date', 'QUIK', 80), ('bond.maturity_date', 'TINVEST', 70),
    ('bond.coupon_rate', 'MOEX', 100), ('bond.coupon_rate', 'BCS', 90),
    ('bond.coupons_per_year', 'MOEX', 100), ('bond.coupons_per_year', 'BCS', 90), ('bond.coupons_per_year', 'TINVEST', 80),
    ('bond.next_coupon_date', 'MOEX', 100), ('bond.next_coupon_date', 'BCS', 90),
    ('bond.coupon_duration', 'QUIK', 100),
    ('derivative.strike', 'MOEX', 100), ('derivative.strike', 'BCS', 90),
    ('derivative.underlying_code', 'MOEX', 100), ('derivative.underlying_code', 'BCS', 90),
    ('future.contract_name', 'MOEX', 100), ('future.start_date', 'MOEX', 100), ('future.expiration_date', 'MOEX', 100), ('future.initial_value', 'MOEX', 100),
    ('future.lot_size', 'BCS', 90),
    ('option.option_type', 'MOEX', 100), ('option.exercise_type', 'MOEX', 100), ('option.margin_style', 'MOEX', 100), ('option.contract_name', 'MOEX', 100), ('option.series_name', 'MOEX', 100),
    ('option.calculation_mode', 'MOEX', 100), ('option.calculation_description', 'MOEX', 100), ('option.settlement_type', 'MOEX', 100), ('option.expiration_date', 'MOEX', 100),
    ('fund.*', 'TINVEST', 100)
) p(attribute_code, source_code, priority)
JOIN ref.external_systems es ON es.ext_system = p.source_code;
CREATE TABLE ref.instrument_source_classifications (
    ext_system_id TINYINT NOT NULL,
    source_code NVARCHAR(128) NOT NULL,
    source_subtype NVARCHAR(128) NOT NULL CONSTRAINT DF_ref_instrument_source_classifications_subtype DEFAULT N'',
    type_id TINYINT NOT NULL,
    subtype_id SMALLINT NULL,
    CONSTRAINT PK_ref_instrument_source_classifications PRIMARY KEY (ext_system_id, source_code, source_subtype),
    CONSTRAINT FK_ref_instrument_source_classifications_system FOREIGN KEY (ext_system_id) REFERENCES ref.external_systems(ext_system_id),
    CONSTRAINT FK_ref_instrument_source_classifications_type FOREIGN KEY (type_id) REFERENCES ref.instrument_types(type_id),
    CONSTRAINT FK_ref_instrument_source_classifications_subtype FOREIGN KEY (type_id, subtype_id) REFERENCES ref.instrument_subtypes(type_id, subtype_id)
);
CREATE TABLE ref.instrument_attribute_sources (
    instrument_id BIGINT NOT NULL,
    attribute_code VARCHAR(80) NOT NULL,
    ext_system_id TINYINT NOT NULL,
    updated_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_instrument_attribute_sources_updated DEFAULT SYSUTCDATETIME(),
    CONSTRAINT PK_ref_instrument_attribute_sources PRIMARY KEY (instrument_id, attribute_code),
    CONSTRAINT FK_ref_instrument_attribute_sources_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments(instrument_id),
    CONSTRAINT FK_ref_instrument_attribute_sources_system FOREIGN KEY (ext_system_id) REFERENCES ref.external_systems(ext_system_id)
);

CREATE TABLE ref.instrument_normalization_conflicts (
    conflict_id BIGINT IDENTITY(1,1) NOT NULL,
    ext_system_id TINYINT NOT NULL,
    external_code NVARCHAR(256) NOT NULL,
    external_code_type_id TINYINT NULL,
    conflict_kind VARCHAR(40) NOT NULL,
    current_instrument_id BIGINT NULL,
    proposed_instrument_id BIGINT NULL,
    attribute_code VARCHAR(80) NULL,
    current_value NVARCHAR(1000) NULL,
    proposed_value NVARCHAR(1000) NULL,
    reason NVARCHAR(1000) NOT NULL,
    status VARCHAR(20) NOT NULL CONSTRAINT DF_ref_instrument_normalization_conflicts_status DEFAULT 'pending',
    first_seen_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_instrument_normalization_conflicts_first DEFAULT SYSUTCDATETIME(),
    last_seen_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_instrument_normalization_conflicts_last DEFAULT SYSUTCDATETIME(),
    conflict_key VARCHAR(64) NOT NULL,
    CONSTRAINT PK_ref_instrument_normalization_conflicts PRIMARY KEY (conflict_id),
    CONSTRAINT FK_ref_instrument_normalization_conflicts_system FOREIGN KEY (ext_system_id) REFERENCES ref.external_systems(ext_system_id),
    CONSTRAINT FK_ref_instrument_normalization_conflicts_current FOREIGN KEY (current_instrument_id) REFERENCES ref.instruments(instrument_id),
    CONSTRAINT FK_ref_instrument_normalization_conflicts_proposed FOREIGN KEY (proposed_instrument_id) REFERENCES ref.instruments(instrument_id),
    CONSTRAINT CK_ref_instrument_normalization_conflicts_status CHECK (status IN ('pending','resolved','ignored'))
);
CREATE INDEX IX_ref_instrument_normalization_conflicts_pending
    ON ref.instrument_normalization_conflicts (status, last_seen_at);
CREATE UNIQUE INDEX UQ_ref_instrument_normalization_conflicts_pending_key
    ON ref.instrument_normalization_conflicts (conflict_key) WHERE status='pending';

CREATE TABLE ref.external_code_claims (
    namespace_code NVARCHAR(50) NOT NULL,
    ext_code_type_id TINYINT NOT NULL,
    ext_code NVARCHAR(100) NOT NULL,
    claimant_system_id TINYINT NOT NULL,
    candidate_instrument_id BIGINT NOT NULL,
    first_seen_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_external_code_claims_first DEFAULT SYSUTCDATETIME(),
    last_seen_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_external_code_claims_last DEFAULT SYSUTCDATETIME(),
    CONSTRAINT PK_ref_external_code_claims PRIMARY KEY (namespace_code, ext_code_type_id, ext_code, claimant_system_id, candidate_instrument_id),
    CONSTRAINT FK_ref_external_code_claims_claimant FOREIGN KEY (claimant_system_id) REFERENCES ref.external_systems(ext_system_id),
    CONSTRAINT FK_ref_external_code_claims_type FOREIGN KEY (ext_code_type_id) REFERENCES ref.external_code_types(ext_code_type_id),
    CONSTRAINT FK_ref_external_code_claims_instrument FOREIGN KEY (candidate_instrument_id) REFERENCES ref.instruments(instrument_id)
);

INSERT INTO ref.external_code_claims (namespace_code, ext_code_type_id, ext_code, claimant_system_id, candidate_instrument_id)
SELECT CASE ec.ext_code_type_id WHEN 4 THEN N'FIGI' WHEN 5 THEN N'ISIN' ELSE CONCAT(N'SYSTEM:',es.ext_system) END,
       ec.ext_code_type_id, ec.ext_code, ec.ext_system_id, ec.internal_id
FROM ref.external_codes ec
JOIN ref.external_systems es ON es.ext_system_id=ec.ext_system_id
WHERE ec.ext_code_type_id IN (3, 4, 5, 6);

CREATE TABLE ref.board_attribute_sources (
    board_id SMALLINT NOT NULL,
    attribute_code VARCHAR(80) NOT NULL,
    ext_system_id TINYINT NOT NULL,
    updated_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_board_attribute_sources_updated DEFAULT SYSUTCDATETIME(),
    CONSTRAINT PK_ref_board_attribute_sources PRIMARY KEY (board_id, attribute_code),
    CONSTRAINT FK_ref_board_attribute_sources_board FOREIGN KEY (board_id) REFERENCES ref.boards(board_id),
    CONSTRAINT FK_ref_board_attribute_sources_system FOREIGN KEY (ext_system_id) REFERENCES ref.external_systems(ext_system_id)
);

CREATE TABLE ref.instrument_board_attribute_sources (
    instrument_id BIGINT NOT NULL,
    board_id SMALLINT NOT NULL,
    attribute_code VARCHAR(80) NOT NULL,
    ext_system_id TINYINT NOT NULL,
    updated_at DATETIME2(3) NOT NULL CONSTRAINT DF_ref_instrument_board_attribute_sources_updated DEFAULT SYSUTCDATETIME(),
    CONSTRAINT PK_ref_instrument_board_attribute_sources PRIMARY KEY (instrument_id, board_id, attribute_code),
    CONSTRAINT FK_ref_instrument_board_attribute_sources_link FOREIGN KEY (instrument_id, board_id) REFERENCES ref.instrument_boards(instrument_id, board_id),
    CONSTRAINT FK_ref_instrument_board_attribute_sources_system FOREIGN KEY (ext_system_id) REFERENCES ref.external_systems(ext_system_id)
);
GRANT SELECT, INSERT, UPDATE ON ref.instruments TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_boards TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.external_codes TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.boards TO quik_instruments_worker;
GRANT SELECT, INSERT ON ref.instrument_type_boards TO quik_instruments_worker;
GRANT SELECT ON ref.instrument_types TO quik_instruments_worker;
GRANT SELECT ON ref.external_systems TO quik_instruments_worker;
GRANT SELECT ON ref.currencies TO quik_instruments_worker;
GRANT SELECT ON ref.instrument_subtypes TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_source_classifications TO quik_instruments_worker;
GRANT SELECT ON ref.instrument_attribute_priorities TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_attribute_sources TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_normalization_conflicts TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.external_code_claims TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.board_attribute_sources TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_board_attribute_sources TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_bond_details TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_fund_details TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_future_details TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_option_details TO quik_instruments_worker;
GRANT SELECT ON SCHEMA::quik TO quik_instruments_worker;
GRANT SELECT ON SCHEMA::bcs TO quik_instruments_worker;
GRANT SELECT ON SCHEMA::tinvest TO quik_instruments_worker;
GRANT SELECT ON SCHEMA::moex TO quik_instruments_worker;

-- +goose Down

THROW 51031, N'Откат миграции 091 остановлен: board_id расширен до SMALLINT с сохранением внешних ссылок; для отката требуется отдельная проверенная миграция', 1;
RETURN;

IF EXISTS (SELECT 1 FROM ref.instrument_bond_details WHERE coupon_rate IS NOT NULL OR coupons_per_year IS NOT NULL OR coupon_type IS NOT NULL OR next_coupon_date IS NOT NULL OR issue_date IS NOT NULL OR nominal_currency IS NOT NULL) THROW 51020, N'Откат невозможен: в деталях облигаций есть значения без legacy-полей', 1;
IF EXISTS (SELECT 1 FROM ref.instrument_bond_details d JOIN ref.instruments i ON i.instrument_id=d.instrument_id WHERE ISNULL(d.face_value,-1)<>ISNULL(i.face_value,-1) OR ISNULL(d.maturity_date,'19000101')<>ISNULL(i.maturity_date,'19000101') OR ISNULL(d.coupon_duration,-1)<>ISNULL(i.coupon_duration,-1)) THROW 51022, N'Откат невозможен: значения деталей облигаций отличаются от legacy-полей', 1;
IF EXISTS (SELECT 1 FROM ref.instrument_fund_details) OR EXISTS (SELECT 1 FROM ref.instrument_future_details) OR EXISTS (SELECT 1 FROM ref.instrument_option_details) THROW 51021, N'Откат невозможен: сначала перенесите типоспецифичные детали', 1;
DROP INDEX IX_ref_instrument_normalization_conflicts_pending ON ref.instrument_normalization_conflicts;
DROP TABLE ref.instrument_normalization_conflicts;
DROP TABLE ref.instrument_attribute_sources;
DROP TABLE ref.instrument_attribute_priorities;
DELETE FROM ref.instrument_bond_details WHERE instrument_id IN (SELECT instrument_id FROM ref.instruments);
DROP TABLE ref.instrument_option_details;
DROP TABLE ref.instrument_future_details;
DROP TABLE ref.instrument_fund_details;
DROP TABLE ref.instrument_bond_details;
ALTER TABLE ref.instruments DROP CONSTRAINT FK_ref_instruments_subtype;
ALTER TABLE ref.instruments DROP COLUMN subtype_id;
DROP TABLE ref.instrument_subtypes;

