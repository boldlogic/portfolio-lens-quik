-- +goose Up

IF SCHEMA_ID(N'bcs') IS NULL
    EXEC(N'CREATE SCHEMA bcs');

IF NOT EXISTS (SELECT 1 FROM ref.external_systems WHERE ext_system = N'BCS')
    INSERT INTO ref.external_systems (ext_system) VALUES (N'BCS');

CREATE TABLE bcs.securities (
    security_id            BIGINT         NOT NULL IDENTITY(1, 1),
    ticker                 NVARCHAR(64)   NOT NULL,
    isin                   NVARCHAR(16)   NULL,
    bcs_type               NVARCHAR(64)   NULL,
    card_subtype           NVARCHAR(64)   NULL,
    quik_subtype           NVARCHAR(60)   NULL,
    ref_instrument_id      BIGINT         NULL,
    display_name           NVARCHAR(200)  NULL,
    short_name             NVARCHAR(200)  NULL,
    display_name_second    NVARCHAR(200)  NULL,
    sub_title              NVARCHAR(200)  NULL,
    registration_code      NVARCHAR(64)   NULL,
    nrd_code               NVARCHAR(64)   NULL,
    cfi                    NVARCHAR(16)   NULL,
    type_name              NVARCHAR(64)   NULL,
    issuer_name            NVARCHAR(200)  NULL,
    business_sector        NVARCHAR(200)  NULL,
    business_sector_id     DECIMAL(19, 8) NULL,
    business_country       NVARCHAR(100)  NULL,
    business_country_code  NVARCHAR(8)    NULL,
    trading_currency       NVARCHAR(16)   NULL,
    currency_nominal       NVARCHAR(16)   NULL,
    currency_step_price    NVARCHAR(16)   NULL,
    settlement_currency    NVARCHAR(16)   NULL,
    first_curr_code        NVARCHAR(16)   NULL,
    primary_board          NVARCHAR(32)   NULL,
    scale                  INT            NULL,
    minimum_step           DECIMAL(19, 8) NULL,
    face_value             DECIMAL(19, 8) NULL,
    lot_size               DECIMAL(19, 8) NULL,
    step_price             DECIMAL(19, 8) NULL,
    CONSTRAINT PK_bcs_securities PRIMARY KEY CLUSTERED (security_id),
    CONSTRAINT FK_bcs_securities_instrument FOREIGN KEY (ref_instrument_id)
        REFERENCES ref.instruments (instrument_id)
);

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_securities_ticker ON bcs.securities (ticker);

CREATE TABLE bcs.security_access (
    security_id                 BIGINT         NOT NULL,
    is_qualified_only           BIT            NULL,
    available_for_unqualified   BIT            NULL,
    is_can_short                BIT            NULL,
    is_can_margin               BIT            NULL,
    is_blocked                  BIT            NULL,
    is_bcs_product              BIT            NULL,
    qualified_test_id           DECIMAL(19, 8) NULL,
    qualified_test_id_tm        DECIMAL(19, 8) NULL,
    complex_product             DECIMAL(19, 8) NULL,
    exclude_type_flags          DECIMAL(19, 8) NULL,
    CONSTRAINT PK_bcs_security_access PRIMARY KEY CLUSTERED (security_id),
    CONSTRAINT FK_bcs_security_access_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.security_bonds (
    security_id           BIGINT         NOT NULL,
    coupon_rate           DECIMAL(19, 8) NULL,
    coupons_per_year      DECIMAL(19, 8) NULL,
    coupon_type_name      NVARCHAR(100)  NULL,
    next_coupon           DATE           NULL,
    accrued_int           DECIMAL(19, 8) NULL,
    maturity_date         DATE           NULL,
    emission_date         DATE           NULL,
    settlement_date       DATE           NULL,
    settle_code           NVARCHAR(16)   NULL,
    credit_rating         NVARCHAR(32)   NULL,
    liquidity_rating      NVARCHAR(32)   NULL,
    is_replacement_bond   BIT            NULL,
    amortised_mty         BIT            NULL,
    CONSTRAINT PK_bcs_security_bonds PRIMARY KEY CLUSTERED (security_id),
    CONSTRAINT FK_bcs_security_bonds_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.security_estimates (
    security_id              BIGINT         NOT NULL,
    bcs_score                DECIMAL(19, 8) NULL,
    bcs_score_color          NVARCHAR(32)   NULL,
    target_price             DECIMAL(19, 8) NULL,
    percent_target_current   DECIMAL(19, 8) NULL,
    promo_idx                DECIMAL(19, 8) NULL,
    dividend_yield           DECIMAL(19, 8) NULL,
    predicted_dps            DECIMAL(19, 8) NULL,
    pe_norm                  DECIMAL(19, 8) NULL,
    mktcap                   DECIMAL(19, 8) NULL,
    price_tangible           DECIMAL(19, 8) NULL,
    eps_growth_rate          DECIMAL(19, 8) NULL,
    price_change_year        DECIMAL(19, 8) NULL,
    price_change_half_year   DECIMAL(19, 8) NULL,
    price_change_month       DECIMAL(19, 8) NULL,
    price_change_early_year  DECIMAL(19, 8) NULL,
    CONSTRAINT PK_bcs_security_estimates PRIMARY KEY CLUSTERED (security_id),
    CONSTRAINT FK_bcs_security_estimates_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.security_derivatives (
    security_id                      BIGINT         NOT NULL,
    base_asset                       NVARCHAR(64)   NULL,
    base_asset_future                NVARCHAR(64)   NULL,
    base_asset_security_sec_code     NVARCHAR(32)   NULL,
    base_asset_security_class_code   NVARCHAR(32)   NULL,
    strike                           DECIMAL(19, 8) NULL,
    CONSTRAINT PK_bcs_security_derivatives PRIMARY KEY CLUSTERED (security_id),
    CONSTRAINT FK_bcs_security_derivatives_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.security_boards (
    security_board_id BIGINT        NOT NULL IDENTITY(1, 1),
    security_id       BIGINT        NOT NULL,
    class_code        NVARCHAR(32)  NOT NULL,
    exchange_code     NVARCHAR(32)  NULL,
    is_secondary      BIT           NOT NULL,
    CONSTRAINT PK_bcs_security_boards PRIMARY KEY CLUSTERED (security_board_id),
    CONSTRAINT FK_bcs_security_boards_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id),
    CONSTRAINT UQ_bcs_security_boards UNIQUE (security_id, class_code)
);

CREATE TABLE bcs.security_exclude_types (
    exclude_id    BIGINT         NOT NULL IDENTITY(1, 1),
    security_id   BIGINT         NOT NULL,
    ord           INT            NOT NULL,
    exclude_type  DECIMAL(19, 8) NOT NULL,
    CONSTRAINT PK_bcs_security_exclude_types PRIMARY KEY CLUSTERED (exclude_id),
    CONSTRAINT FK_bcs_security_exclude_types_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id),
    CONSTRAINT UQ_bcs_security_exclude_types UNIQUE (security_id, ord)
);

CREATE TABLE bcs.portfolio (
    portfolio_id BIGINT            NOT NULL,
    received_at  DATETIMEOFFSET(7) NOT NULL,
    http_status  INT               NOT NULL,
    raw_body     VARBINARY(MAX)    NOT NULL,
    CONSTRAINT PK_bcs_portfolio PRIMARY KEY CLUSTERED (portfolio_id)
);

CREATE TABLE bcs.positions (
    position_id            BIGINT         NOT NULL IDENTITY(1, 1),
    security_id            BIGINT         NOT NULL,
    position_type          NVARCHAR(32)   NOT NULL,
    account_code           NVARCHAR(32)   NULL,
    exchange_code          NVARCHAR(32)   NULL,
    board_code             NVARCHAR(32)   NULL,
    currency_code          NVARCHAR(16)   NULL,
    term                   NVARCHAR(8)    NOT NULL,
    quantity               DECIMAL(19, 8) NOT NULL,
    current_value_rub      DECIMAL(19, 8) NOT NULL,
    unrealized_pl          DECIMAL(19, 8) NOT NULL,
    balance_price          DECIMAL(19, 8) NULL,
    current_price          DECIMAL(19, 8) NULL,
    scale                  INT            NULL,
    balance_value_rub      DECIMAL(19, 8) NULL,
    daily_pl               DECIMAL(19, 8) NULL,
    unrealized_percent_pl  DECIMAL(19, 8) NULL,
    accrued_income         DECIMAL(19, 8) NULL,
    locked                 DECIMAL(19, 8) NULL,
    CONSTRAINT PK_bcs_positions PRIMARY KEY CLUSTERED (position_id),
    CONSTRAINT FK_bcs_positions_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.limit_payload (
    payload_id BIGINT         NOT NULL,
    raw_body   VARBINARY(MAX) NOT NULL,
    CONSTRAINT PK_bcs_limit_payload PRIMARY KEY CLUSTERED (payload_id)
);

CREATE TABLE bcs.depo_limits (
    depo_limit_id        BIGINT            NOT NULL IDENTITY(1, 1),
    security_id          BIGINT            NOT NULL,
    class_code           NVARCHAR(32)      NULL,
    exchange_code        NVARCHAR(32)      NULL,
    average_price        DECIMAL(19, 8)    NULL,
    scale                INT               NULL,
    quantity_type        NVARCHAR(32)      NOT NULL,
    quantity_value       DECIMAL(19, 8)    NULL,
    batch_type           NVARCHAR(32)      NOT NULL,
    batch_value          DECIMAL(19, 8)    NULL,
    instrument_type      NVARCHAR(64)      NOT NULL,
    load_date            DATETIMEOFFSET(7) NULL,
    locked_buy_value     DECIMAL(19, 8)    NULL,
    locked_sell_value    DECIMAL(19, 8)    NULL,
    locked_buy_quantity  DECIMAL(19, 8)    NULL,
    locked_sell_quantity DECIMAL(19, 8)    NULL,
    CONSTRAINT PK_bcs_depo_limits PRIMARY KEY CLUSTERED (depo_limit_id),
    CONSTRAINT FK_bcs_depo_limits_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.money_limits (
    money_limit_id  BIGINT            NOT NULL IDENTITY(1, 1),
    exchange_code   NVARCHAR(32)      NULL,
    currency_code   NVARCHAR(16)      NULL,
    locked          DECIMAL(19, 8)    NULL,
    average_price   DECIMAL(19, 8)    NULL,
    instrument_type NVARCHAR(64)      NOT NULL,
    quantity_type   NVARCHAR(32)      NOT NULL,
    quantity_value  DECIMAL(19, 8)    NULL,
    load_date       DATETIMEOFFSET(7) NULL,
    CONSTRAINT PK_bcs_money_limits PRIMARY KEY CLUSTERED (money_limit_id)
);

CREATE TABLE bcs.future_holdings (
    holding_id       BIGINT            NOT NULL IDENTITY(1, 1),
    security_id      BIGINT            NOT NULL,
    class_code       NVARCHAR(32)      NULL,
    exchange_code    NVARCHAR(32)      NULL,
    average_price    DECIMAL(19, 8)    NULL,
    scale            INT               NULL,
    total_net        DECIMAL(19, 8)    NULL,
    var_margin       DECIMAL(19, 8)    NULL,
    position_value   DECIMAL(19, 8)    NULL,
    instrument_type  NVARCHAR(64)      NOT NULL,
    trade_date       DATETIMEOFFSET(7) NULL,
    execution_date   DATETIMEOFFSET(7) NULL,
    CONSTRAINT PK_bcs_future_holdings PRIMARY KEY CLUSTERED (holding_id),
    CONSTRAINT FK_bcs_future_holdings_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id)
);

CREATE TABLE bcs.futures_limits (
    futures_limit_id BIGINT            NOT NULL IDENTITY(1, 1),
    currency_code    NVARCHAR(16)      NULL,
    exchange_code    NVARCHAR(32)      NULL,
    instrument_type  NVARCHAR(64)      NOT NULL,
    cbp_limit        DECIMAL(19, 8)    NULL,
    var_margin       DECIMAL(19, 8)    NULL,
    load_date        DATETIMEOFFSET(7) NULL,
    CONSTRAINT PK_bcs_futures_limits PRIMARY KEY CLUSTERED (futures_limit_id)
);

CREATE TABLE bcs.quotes (
    quote_id      BIGINT            NOT NULL IDENTITY(1, 1),
    security_id   BIGINT            NOT NULL,
    class_code    NVARCHAR(32)      NOT NULL,
    quote_day     DATE              NOT NULL,
    scale         INT               NULL,
    last_price    DECIMAL(19, 8)    NULL,
    bid           DECIMAL(19, 8)    NULL,
    offer         DECIMAL(19, 8)    NULL,
    change_amount DECIMAL(19, 8)    NULL,
    quoted_at     DATETIMEOFFSET(7) NOT NULL,
    CONSTRAINT PK_bcs_quotes PRIMARY KEY CLUSTERED (quote_id),
    CONSTRAINT FK_bcs_quotes_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id),
    CONSTRAINT UQ_bcs_quotes_day UNIQUE (security_id, class_code, quote_day)
);

CREATE TABLE bcs.schedules (
    schedule_id  BIGINT       NOT NULL IDENTITY(1, 1),
    security_id  BIGINT       NOT NULL,
    class_code   NVARCHAR(32) NOT NULL,
    is_work_day  BIT          NULL,
    CONSTRAINT PK_bcs_schedules PRIMARY KEY CLUSTERED (schedule_id),
    CONSTRAINT FK_bcs_schedules_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id),
    CONSTRAINT UQ_bcs_schedules UNIQUE (security_id, class_code)
);

CREATE TABLE bcs.schedule_sessions (
    session_id     BIGINT       NOT NULL IDENTITY(1, 1),
    schedule_id    BIGINT       NOT NULL,
    session_ord    INT          NOT NULL,
    start_time     NVARCHAR(16) NOT NULL,
    end_time       NVARCHAR(16) NOT NULL,
    session_type   NVARCHAR(80) NOT NULL,
    session_status NVARCHAR(16) NOT NULL,
    CONSTRAINT PK_bcs_schedule_sessions PRIMARY KEY CLUSTERED (session_id),
    CONSTRAINT FK_bcs_schedule_sessions_schedule FOREIGN KEY (schedule_id)
        REFERENCES bcs.schedules (schedule_id),
    CONSTRAINT UQ_bcs_schedule_sessions UNIQUE (schedule_id, session_ord)
);

CREATE TABLE bcs.class_trading_status (
    status_id        BIGINT            NOT NULL IDENTITY(1, 1),
    class_code       NVARCHAR(32)      NOT NULL,
    session_type_id  INT               NULL,
    session_type     NVARCHAR(80)      NOT NULL,
    session_status   NVARCHAR(16)      NOT NULL,
    next_session_at  DATETIMEOFFSET(7) NULL,
    CONSTRAINT PK_bcs_class_trading_status PRIMARY KEY CLUSTERED (status_id),
    CONSTRAINT UQ_bcs_class_trading_status_class UNIQUE (class_code)
);

CREATE TABLE bcs.own_trades (
    trade_id             BIGINT            NOT NULL IDENTITY(1, 1),
    trade_num            BIGINT            NOT NULL,
    order_num            BIGINT            NULL,
    security_id          BIGINT            NULL,
    ticker               NVARCHAR(64)      NOT NULL,
    client_code          NVARCHAR(32)      NOT NULL,
    class_code           NVARCHAR(32)      NULL,
    settlement_currency  NVARCHAR(16)      NULL,
    base_currency        NVARCHAR(16)      NULL,
    price_currency       NVARCHAR(16)      NULL,
    side                 NVARCHAR(8)       NOT NULL,
    instrument_type      NVARCHAR(64)      NOT NULL,
    deal_type            INT               NULL,
    trade_date_time      DATETIMEOFFSET(7) NOT NULL,
    price                DECIMAL(19, 8)    NULL,
    volume               DECIMAL(19, 8)    NULL,
    go_amount            DECIMAL(19, 8)    NULL,
    contract_amount      DECIMAL(19, 8)    NULL,
    settle_date          DATE              NULL,
    trade_quantity       DECIMAL(19, 8)    NULL,
    trade_quantity_lots  DECIMAL(19, 8)    NULL,
    CONSTRAINT PK_bcs_own_trades PRIMARY KEY CLUSTERED (trade_id),
    CONSTRAINT FK_bcs_own_trades_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id),
    CONSTRAINT UQ_bcs_own_trades_num UNIQUE (trade_num)
);

CREATE TABLE bcs.cash_operations (
    operation_id        BIGINT            NOT NULL IDENTITY(1, 1),
    broker_id           NVARCHAR(64)      NOT NULL,
    operation_date      DATETIMEOFFSET(7) NOT NULL,
    is_iia              BIT               NULL,
    security_id         BIGINT            NULL,
    ticker              NVARCHAR(64)      NULL,
    class_code          NVARCHAR(32)      NULL,
    operation_type      NVARCHAR(64)      NOT NULL,
    status              NVARCHAR(32)      NOT NULL,
    sum_amount          DECIMAL(19, 8)    NULL,
    currency_code       NVARCHAR(16)      NULL,
    is_rests_pay_out    BIT               NULL,
    isin                NVARCHAR(16)      NULL,
    issuer_name         NVARCHAR(200)     NULL,
    replenishment_type  NVARCHAR(64)      NULL,
    execution_period    NVARCHAR(64)      NULL,
    balance_change      NVARCHAR(16)      NOT NULL,
    rejection_reason    NVARCHAR(400)     NULL,
    CONSTRAINT PK_bcs_cash_operations PRIMARY KEY CLUSTERED (operation_id),
    CONSTRAINT FK_bcs_cash_operations_security FOREIGN KEY (security_id)
        REFERENCES bcs.securities (security_id),
    CONSTRAINT UQ_bcs_cash_operations_broker UNIQUE (broker_id)
);

CREATE TABLE bcs.import_state (
    import_name NVARCHAR(64)      NOT NULL,
    imported_at DATETIMEOFFSET(7) NOT NULL,
    CONSTRAINT PK_bcs_import_state PRIMARY KEY CLUSTERED (import_name)
);

-- +goose Down

DROP TABLE IF EXISTS bcs.import_state;
DROP TABLE IF EXISTS bcs.cash_operations;
DROP TABLE IF EXISTS bcs.own_trades;
DROP TABLE IF EXISTS bcs.class_trading_status;
DROP TABLE IF EXISTS bcs.schedule_sessions;
DROP TABLE IF EXISTS bcs.schedules;
DROP TABLE IF EXISTS bcs.quotes;
DROP TABLE IF EXISTS bcs.futures_limits;
DROP TABLE IF EXISTS bcs.future_holdings;
DROP TABLE IF EXISTS bcs.money_limits;
DROP TABLE IF EXISTS bcs.depo_limits;
DROP TABLE IF EXISTS bcs.limit_payload;
DROP TABLE IF EXISTS bcs.positions;
DROP TABLE IF EXISTS bcs.portfolio;
DROP TABLE IF EXISTS bcs.security_exclude_types;
DROP TABLE IF EXISTS bcs.security_boards;
DROP TABLE IF EXISTS bcs.security_derivatives;
DROP TABLE IF EXISTS bcs.security_estimates;
DROP TABLE IF EXISTS bcs.security_bonds;
DROP TABLE IF EXISTS bcs.security_access;
DROP TABLE IF EXISTS bcs.securities;

IF SCHEMA_ID(N'bcs') IS NOT NULL
    EXEC(N'DROP SCHEMA bcs');

DECLARE @bcs TINYINT = (SELECT ext_system_id FROM ref.external_systems WHERE ext_system = N'BCS');
IF @bcs IS NOT NULL AND NOT EXISTS (SELECT 1 FROM ref.external_codes WHERE ext_system_id = @bcs)
    DELETE FROM ref.external_systems WHERE ext_system_id = @bcs;
