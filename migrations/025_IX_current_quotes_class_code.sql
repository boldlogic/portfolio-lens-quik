-- +goose Up
CREATE TABLE
    ref.instruments (
        instrument_id BIGINT IDENTITY (1, 1),
        trade_point_id tinyint NOT NULL,
        ticker NVARCHAR (12) NOT NULL, -- Код инструмента (char 15)
        registration_number NVARCHAR (30) NULL, -- Рег.номер (char 250)
        full_name NVARCHAR (100) NULL, -- Инструмент (char 250)
        short_name NVARCHAR (50) NULL, -- Инструмент сокр. (char 100)
        isin NVARCHAR (12) NULL, -- ISIN (char 15)
        face_value decimal(19, 8) NULL, -- Номинал
        maturity_date DATE NULL, -- Погашение (date)
        coupon_duration INT NULL, -- Длит. купона (int)
        rw ROWVERSION NOT NULL,
        CONSTRAINT PK_ref_instruments PRIMARY KEY CLUSTERED (instrument_id),
        CONSTRAINT FK_ref_instruments_trade_point_id FOREIGN KEY (trade_point_id) REFERENCES ref.trade_points (point_id),
    );

CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_instruments_ticker_trade_point_id ON ref.instruments (ticker, trade_point_id);

CREATE TABLE
    ref.instrument_boards (
        instrument_id BIGINT NOT NULL,
        board_id tinyint NOT NULL,
        currency_id SMALLINT NULL,
        base_currency_id SMALLINT NULL,
        quote_currency_id SMALLINT NULL,
        counter_currency_id SMALLINT NULL,
        is_primary bit NULL,
        CONSTRAINT PK_ref_instrument_boards PRIMARY KEY CLUSTERED (instrument_id, board_id),
        CONSTRAINT FK_ref_instrument_boards_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments (instrument_id),
        CONSTRAINT FK_ref_instrument_boards_board FOREIGN KEY (board_id) REFERENCES ref.boards (board_id),
        CONSTRAINT FK_ref_instrument_boards_currency FOREIGN KEY (currency_id) REFERENCES ref.currencies (iso_code),
        CONSTRAINT FK_ref_instrument_boards_base_currency FOREIGN KEY (base_currency_id) REFERENCES ref.currencies (iso_code),
        CONSTRAINT FK_ref_instrument_boards_quote_currency FOREIGN KEY (quote_currency_id) REFERENCES ref.currencies (iso_code),
        CONSTRAINT FK_ref_instrument_boards_counter_currency FOREIGN KEY (counter_currency_id) REFERENCES ref.currencies (iso_code),
    );