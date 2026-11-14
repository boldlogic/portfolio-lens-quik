-- +goose Up

CREATE SCHEMA tinvest;


CREATE TABLE tinvest.account (

    id NVARCHAR(128) PRIMARY KEY,
    name NVARCHAR(200) NOT NULL DEFAULT '',
    type_code NVARCHAR(64) NOT NULL
);

CREATE TABLE tinvest.instrument (

    uid NVARCHAR(64) PRIMARY KEY,
    figi NVARCHAR(32) NOT NULL DEFAULT '',
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    name NVARCHAR(200) NOT NULL DEFAULT '',
    class_code NVARCHAR(32) NOT NULL DEFAULT '',
    kind NVARCHAR(48) NOT NULL DEFAULT ''
);

CREATE TABLE tinvest.bond_passport (

    instrument_uid NVARCHAR(64) PRIMARY KEY REFERENCES tinvest.instrument(uid),
    nominal_amount DECIMAL(28, 9) NOT NULL,
    nominal_currency NVARCHAR(16) NOT NULL,
    maturity DATE NULL,
    coupons_per_year INT NOT NULL
);

CREATE TABLE tinvest.user_profile (

    user_id NVARCHAR(128) PRIMARY KEY,
    premium BIT NOT NULL,
    qualified BIT NOT NULL,
    tariff NVARCHAR(128) NOT NULL DEFAULT ''
);

CREATE TABLE tinvest.user_qualification (

    user_id NVARCHAR(128) NOT NULL REFERENCES tinvest.user_profile(user_id),
    permission NVARCHAR(128) NOT NULL,
    PRIMARY KEY (user_id, permission)
);

CREATE TABLE tinvest.rate_limit (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    user_id NVARCHAR(128) NOT NULL REFERENCES tinvest.user_profile(user_id),
    kind NVARCHAR(48) NOT NULL,
    per_minute BIGINT NOT NULL DEFAULT 0,
    per_second BIGINT NOT NULL DEFAULT 0,
    stream_limit BIGINT NOT NULL DEFAULT 0,
    stream_open BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE tinvest.rate_limit_method (

    limit_id BIGINT NOT NULL REFERENCES tinvest.rate_limit(id),
    method NVARCHAR(128) NOT NULL,
    PRIMARY KEY (limit_id, method)
);

CREATE TABLE tinvest.coupon (

    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    pay_date DATE NOT NULL,
    amount DECIMAL(28, 9) NOT NULL,
    currency NVARCHAR(16) NOT NULL,
    coupon_type NVARCHAR(128) NOT NULL DEFAULT '',
    PRIMARY KEY (instrument_uid, pay_date)
);

CREATE TABLE tinvest.bond_event (

    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    event_type NVARCHAR(128) NOT NULL,
    event_date DATE NOT NULL,
    pay_amount DECIMAL(28, 9) NOT NULL,
    pay_currency NVARCHAR(16) NOT NULL,
    PRIMARY KEY (instrument_uid, event_type, event_date)
);

CREATE TABLE tinvest.dividend (

    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    record_date DATE NOT NULL,
    dividend_type NVARCHAR(128) NOT NULL,
    payment_date DATE NULL,
    amount DECIMAL(28, 9) NOT NULL,
    currency NVARCHAR(16) NOT NULL,
    PRIMARY KEY (instrument_uid, record_date, dividend_type)
);

CREATE TABLE tinvest.snapshot (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    checked_at DATETIMEOFFSET(7) NOT NULL,
    endpoint NVARCHAR(128) NOT NULL DEFAULT ''
);

CREATE TABLE tinvest.snapshot_account (

    snapshot_id BIGINT NOT NULL REFERENCES tinvest.snapshot(id),
    account_id NVARCHAR(128) NOT NULL REFERENCES tinvest.account(id),
    skipped BIT NOT NULL,
    skip_reason NVARCHAR(128) NOT NULL DEFAULT '',
    limits_loading BIT NOT NULL,
    value_check_partial BIT NOT NULL,
    PRIMARY KEY (snapshot_id, account_id)
);

CREATE TABLE tinvest.money_amount (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    bucket NVARCHAR(48) NOT NULL,
    currency NVARCHAR(16) NOT NULL,
    amount DECIMAL(28, 9) NOT NULL,
    PRIMARY KEY (snapshot_id, account_id, bucket, currency),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.percent_reading (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    kind NVARCHAR(48) NOT NULL,
    value DECIMAL(28, 9) NOT NULL,
    PRIMARY KEY (snapshot_id, account_id, kind),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.position (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    quantity DECIMAL(28, 9) NOT NULL,
    current_amount DECIMAL(28, 9) NOT NULL,
    current_currency NVARCHAR(16) NOT NULL,
    nkd_amount DECIMAL(28, 9) NOT NULL,
    nkd_currency NVARCHAR(16) NOT NULL,
    market_amount DECIMAL(28, 9) NOT NULL,
    market_currency NVARCHAR(16) NOT NULL,
    market_known BIT NOT NULL,
    last_price DECIMAL(28, 9) NULL,
    last_known BIT NOT NULL,
    close_price DECIMAL(28, 9) NULL,
    close_known BIT NOT NULL,
    expected_yield DECIMAL(28, 9) NULL,
    expected_known BIT NOT NULL,
    blocked BIT NOT NULL,
    cash BIT NOT NULL,
    trading_status NVARCHAR(128) NOT NULL DEFAULT '',
    PRIMARY KEY (snapshot_id, account_id, instrument_uid),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.holding (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    kind NVARCHAR(48) NOT NULL,
    balance BIGINT NOT NULL,
    blocked BIT NOT NULL,
    exchange_blocked BIT NOT NULL,
    PRIMARY KEY (snapshot_id, account_id, instrument_uid, kind),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.virtual_position (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    quantity DECIMAL(28, 9) NOT NULL,
    current_amount DECIMAL(28, 9) NOT NULL,
    current_currency NVARCHAR(16) NOT NULL,
    expire_date DATE NULL,
    PRIMARY KEY (snapshot_id, account_id, instrument_uid),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.snapshot_order (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    order_id NVARCHAR(64) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT '',
    direction NVARCHAR(128) NOT NULL DEFAULT '',
    lots_requested INT NOT NULL,
    lots_executed INT NOT NULL,
    price_amount DECIMAL(28, 9) NOT NULL,
    price_currency NVARCHAR(16) NOT NULL,
    status NVARCHAR(64) NOT NULL DEFAULT '',
    ordered_at DATETIMEOFFSET(7) NULL,
    PRIMARY KEY (snapshot_id, order_id),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.snapshot_stop_order (

    snapshot_id BIGINT NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    stop_order_id NVARCHAR(64) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT '',
    direction NVARCHAR(128) NOT NULL DEFAULT '',
    lots_requested INT NOT NULL,
    stop_amount NVARCHAR(128) NOT NULL,
    stop_currency NVARCHAR(16) NOT NULL,
    status NVARCHAR(64) NOT NULL DEFAULT '',
    ordered_at DATETIMEOFFSET(7) NULL,
    PRIMARY KEY (snapshot_id, stop_order_id),
    FOREIGN KEY (snapshot_id, account_id) REFERENCES tinvest.snapshot_account(snapshot_id, account_id)
);

CREATE TABLE tinvest.snapshot_warning (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES tinvest.snapshot(id),
    account_id NVARCHAR(128) NOT NULL DEFAULT '',
    message NVARCHAR(2000) NOT NULL
);

CREATE TABLE tinvest.operation (

    account_id NVARCHAR(128) NOT NULL,
    executed_at DATETIMEOFFSET(7) NOT NULL,
    type NVARCHAR(64) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT '',
    payment_amount NVARCHAR(128) NOT NULL,
    payment_currency NVARCHAR(16) NOT NULL,
    quantity DECIMAL(28, 9) NOT NULL,
    name NVARCHAR(200) NOT NULL DEFAULT '',
    state NVARCHAR(64) NOT NULL DEFAULT '',
    figi NVARCHAR(32) NOT NULL DEFAULT '',
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    UNIQUE (account_id, executed_at, type, instrument_uid, payment_amount, payment_currency, quantity)
);

CREATE TABLE tinvest.broker_trade (

    trade_id NVARCHAR(64) NOT NULL DEFAULT '',
    account_id NVARCHAR(128) NOT NULL,
    executed_at DATETIMEOFFSET(7) NOT NULL,
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    direction NVARCHAR(128) NOT NULL DEFAULT '',
    quantity DECIMAL(28, 9) NOT NULL,
    amount DECIMAL(28, 9) NOT NULL,
    currency NVARCHAR(16) NOT NULL,
    broker_commission DECIMAL(28, 9) NOT NULL DEFAULT '0',
    broker_currency NVARCHAR(16) NOT NULL DEFAULT '',
    exchange_commission DECIMAL(28, 9) NOT NULL DEFAULT '0',
    exchange_currency NVARCHAR(16) NOT NULL DEFAULT ''
);

CREATE TABLE tinvest.share_passport (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    instrument_uid NVARCHAR(64) NOT NULL UNIQUE REFERENCES tinvest.instrument(uid),
    lot INT NOT NULL,
    currency NVARCHAR(16) NOT NULL DEFAULT '',
    class_code NVARCHAR(32) NOT NULL DEFAULT '',
    exchange NVARCHAR(32) NOT NULL DEFAULT '',
    isin NVARCHAR(128) NOT NULL DEFAULT '',
    sector NVARCHAR(200) NOT NULL DEFAULT '',
    share_type NVARCHAR(128) NOT NULL DEFAULT '',
    dividend_yield INT NOT NULL,
    nominal_amount DECIMAL(28, 9) NULL,
    nominal_currency NVARCHAR(16) NOT NULL DEFAULT '',
    for_qual BIT NOT NULL
);

CREATE TABLE tinvest.etf_passport (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    instrument_uid NVARCHAR(64) NOT NULL UNIQUE REFERENCES tinvest.instrument(uid),
    lot INT NOT NULL,
    currency NVARCHAR(16) NOT NULL DEFAULT '',
    class_code NVARCHAR(32) NOT NULL DEFAULT '',
    exchange NVARCHAR(32) NOT NULL DEFAULT '',
    isin NVARCHAR(128) NOT NULL DEFAULT '',
    sector NVARCHAR(200) NOT NULL DEFAULT '',
    focus_type NVARCHAR(128) NOT NULL DEFAULT '',
    commission DECIMAL(28, 9) NULL,
    rebalancing NVARCHAR(128) NOT NULL DEFAULT ''
);

CREATE TABLE tinvest.user_risk_rate (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    user_id NVARCHAR(128) NOT NULL REFERENCES tinvest.user_profile(user_id),
    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    side NVARCHAR(128) NOT NULL,
    role NVARCHAR(128) NOT NULL,
    level_code NVARCHAR(128) NOT NULL,
    value DECIMAL(28, 9) NOT NULL,
    UNIQUE (user_id, instrument_uid, side, role, level_code)
);

CREATE TABLE tinvest.trading_day (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    exchange NVARCHAR(32) NOT NULL,
    day DATE NOT NULL,
    is_trading BIT NOT NULL,
    start_at DATETIMEOFFSET(7) NULL,
    end_at DATETIMEOFFSET(7) NULL,
    UNIQUE (exchange, day)
);

CREATE TABLE tinvest.fill (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    trade_num NVARCHAR(64) NOT NULL DEFAULT '',
    account_id NVARCHAR(128) NOT NULL,
    executed_at DATETIMEOFFSET(7) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT '',
    figi NVARCHAR(32) NOT NULL DEFAULT '',
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    direction NVARCHAR(128) NOT NULL DEFAULT '',
    price_amount DECIMAL(28, 9) NOT NULL,
    price_currency NVARCHAR(16) NOT NULL DEFAULT '',
    quantity DECIMAL(28, 9) NOT NULL
);

CREATE TABLE tinvest.stream_event (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    received_at DATETIMEOFFSET(7) NOT NULL,
    stream NVARCHAR(32) NOT NULL,
    account_id NVARCHAR(128) NOT NULL DEFAULT '',
    event_at DATETIMEOFFSET(7) NULL,
    type NVARCHAR(64) NOT NULL DEFAULT '',
    name NVARCHAR(200) NOT NULL DEFAULT '',
    state NVARCHAR(64) NOT NULL DEFAULT '',
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    figi NVARCHAR(32) NOT NULL DEFAULT '',
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT '',
    payment_amount NVARCHAR(128) NOT NULL DEFAULT '',
    payment_currency NVARCHAR(16) NOT NULL DEFAULT '',
    total_amount NVARCHAR(128) NOT NULL DEFAULT '',
    total_currency NVARCHAR(16) NOT NULL DEFAULT '',
    money_count INT NOT NULL DEFAULT 0,
    securities INT NOT NULL DEFAULT 0,
    futures INT NOT NULL DEFAULT 0,
    options INT NOT NULL DEFAULT 0
);

CREATE TABLE tinvest.instrument_asset (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    instrument_uid NVARCHAR(64) NOT NULL UNIQUE REFERENCES tinvest.instrument(uid),
    asset_uid NVARCHAR(64) NOT NULL
);

CREATE TABLE tinvest.accrued_interest (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    accrued_at DATETIMEOFFSET(7) NOT NULL,
    amount DECIMAL(28, 9) NULL,
    value_percent DECIMAL(28, 9) NULL,
    nominal_amount DECIMAL(28, 9) NULL,
    UNIQUE (instrument_uid, accrued_at)
);

CREATE TABLE tinvest.asset_fundamental (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    asset_uid NVARCHAR(64) NOT NULL,
    metric NVARCHAR(128) NOT NULL,
    value_text NVARCHAR(200) NOT NULL,
    UNIQUE (asset_uid, metric)
);

CREATE TABLE tinvest.insider_deal (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    trade_id BIGINT NOT NULL DEFAULT 0,
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT '',
    direction NVARCHAR(128) NOT NULL DEFAULT '',
    currency NVARCHAR(16) NOT NULL DEFAULT '',
    executed_at DATETIMEOFFSET(7) NULL,
    quantity DECIMAL(28, 9) NOT NULL,
    price_amount DECIMAL(28, 9) NOT NULL,
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    investor_name NVARCHAR(200) NOT NULL DEFAULT '',
    investor_position NVARCHAR(200) NOT NULL DEFAULT '',
    percentage DECIMAL(28, 9) NULL,
    option_execution BIT NOT NULL,
    disclosure_at DATETIMEOFFSET(7) NULL
);

CREATE TABLE tinvest.news (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    news_id BIGINT NOT NULL DEFAULT 0,
    source NVARCHAR(128) NOT NULL DEFAULT '',
    title NVARCHAR(400) NOT NULL DEFAULT '',
    content NVARCHAR(MAX) NOT NULL DEFAULT '',
    summary NVARCHAR(1000) NOT NULL DEFAULT '',
    priority BIT NOT NULL,
    published_at DATETIMEOFFSET(7) NULL
);

CREATE TABLE tinvest.news_instrument (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    news_id BIGINT NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL REFERENCES tinvest.instrument(uid),
    UNIQUE (news_id, instrument_uid)
);

CREATE TABLE tinvest.instrument_search (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    query NVARCHAR(128) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL,
    isin NVARCHAR(128) NOT NULL DEFAULT '',
    figi NVARCHAR(32) NOT NULL DEFAULT '',
    ticker NVARCHAR(64) NOT NULL DEFAULT '',
    class_code NVARCHAR(32) NOT NULL DEFAULT '',
    instrument_type NVARCHAR(128) NOT NULL DEFAULT '',
    name NVARCHAR(200) NOT NULL DEFAULT '',
    position_uid NVARCHAR(128) NOT NULL DEFAULT '',
    kind NVARCHAR(48) NOT NULL DEFAULT '',
    api_trade BIT NOT NULL,
    for_iis BIT NOT NULL,
    for_qual BIT NOT NULL,
    weekend BIT NOT NULL,
    blocked_tca BIT NOT NULL,
    lot INT NOT NULL,
    first_1min DATETIMEOFFSET(7) NULL,
    first_1day DATETIMEOFFSET(7) NULL,
    UNIQUE (query, instrument_uid)
);

CREATE TABLE tinvest.country (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    alfa_two NVARCHAR(128) NOT NULL DEFAULT '',
    alfa_three NVARCHAR(128) NOT NULL DEFAULT '',
    name NVARCHAR(200) NOT NULL DEFAULT '',
    name_brief NVARCHAR(128) NOT NULL DEFAULT ''
);

CREATE TABLE tinvest.brand (

    id BIGINT IDENTITY(1, 1) PRIMARY KEY,
    uid NVARCHAR(64) NOT NULL DEFAULT '',
    name NVARCHAR(200) NOT NULL DEFAULT '',
    description NVARCHAR(1000) NOT NULL DEFAULT '',
    info NVARCHAR(2000) NOT NULL DEFAULT '',
    company NVARCHAR(200) NOT NULL DEFAULT '',
    sector NVARCHAR(200) NOT NULL DEFAULT '',
    country_of_risk NVARCHAR(128) NOT NULL DEFAULT '',
    country_of_risk_name NVARCHAR(128) NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX broker_trade_id ON tinvest.broker_trade (trade_id) WHERE trade_id <> '';

CREATE UNIQUE INDEX broker_trade_fallback ON tinvest.broker_trade (account_id, executed_at, ticker, direction, quantity, amount, currency)
    WHERE trade_id = '';

CREATE UNIQUE INDEX fill_num ON tinvest.fill (trade_num) WHERE trade_num <> '';

CREATE UNIQUE INDEX fill_fallback ON tinvest.fill (account_id, executed_at, instrument_uid, price_amount, price_currency, quantity)
    WHERE trade_num = '';

CREATE UNIQUE INDEX insider_trade ON tinvest.insider_deal (trade_id) WHERE trade_id <> 0;

CREATE UNIQUE INDEX insider_fallback ON tinvest.insider_deal (instrument_uid, executed_at, investor_name, direction, quantity, price_amount)
    WHERE trade_id = 0;

CREATE UNIQUE INDEX news_ext ON tinvest.news (news_id) WHERE news_id <> 0;

CREATE UNIQUE INDEX country_alfa_two ON tinvest.country (alfa_two) WHERE alfa_two <> '';

CREATE UNIQUE INDEX country_alfa_three ON tinvest.country (alfa_three) WHERE alfa_two = '' AND alfa_three <> '';

CREATE UNIQUE INDEX brand_uid ON tinvest.brand (uid) WHERE uid <> '';


ALTER TABLE tinvest.account ADD source_account_id BIGINT NULL;
ALTER TABLE tinvest.account ADD CONSTRAINT FK_tinvest_account_source_account FOREIGN KEY (source_account_id)
    REFERENCES portfolio.source_accounts (source_account_id);

ALTER TABLE tinvest.instrument ADD ref_instrument_id BIGINT NULL;
ALTER TABLE tinvest.instrument ADD CONSTRAINT FK_tinvest_instrument_ref FOREIGN KEY (ref_instrument_id)
    REFERENCES ref.instruments (instrument_id);

GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON SCHEMA::tinvest TO quik_portfolio_writer;

-- +goose Down


DROP TABLE tinvest.brand;
DROP TABLE tinvest.country;
DROP TABLE tinvest.instrument_search;
DROP TABLE tinvest.news_instrument;
DROP TABLE tinvest.news;
DROP TABLE tinvest.insider_deal;
DROP TABLE tinvest.asset_fundamental;
DROP TABLE tinvest.accrued_interest;
DROP TABLE tinvest.instrument_asset;
DROP TABLE tinvest.stream_event;
DROP TABLE tinvest.fill;
DROP TABLE tinvest.trading_day;
DROP TABLE tinvest.user_risk_rate;
DROP TABLE tinvest.etf_passport;
DROP TABLE tinvest.share_passport;
DROP TABLE tinvest.broker_trade;
DROP TABLE tinvest.operation;
DROP TABLE tinvest.snapshot_warning;
DROP TABLE tinvest.snapshot_stop_order;
DROP TABLE tinvest.snapshot_order;
DROP TABLE tinvest.virtual_position;
DROP TABLE tinvest.holding;
DROP TABLE tinvest.position;
DROP TABLE tinvest.percent_reading;
DROP TABLE tinvest.money_amount;
DROP TABLE tinvest.snapshot_account;
DROP TABLE tinvest.snapshot;
DROP TABLE tinvest.dividend;
DROP TABLE tinvest.bond_event;
DROP TABLE tinvest.coupon;
DROP TABLE tinvest.rate_limit_method;
DROP TABLE tinvest.rate_limit;
DROP TABLE tinvest.user_qualification;
DROP TABLE tinvest.user_profile;
DROP TABLE tinvest.bond_passport;
DROP TABLE tinvest.instrument;
DROP TABLE tinvest.account;
DROP SCHEMA tinvest;
