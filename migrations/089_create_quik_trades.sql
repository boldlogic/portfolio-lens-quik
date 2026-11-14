-- +goose Up

CREATE TABLE quik.trades (
    trade_num                       decimal(19, 0) NOT NULL,
    exchange_code                   varchar(16)    NULL,
    session_date                    date           NULL,
    trade_date                      date           NULL,
    settle_date                     date           NULL,
    trade_time                      time(0)        NULL,
    trade_time_mcs                  int            NULL,
    order_num                       decimal(19, 0) NULL,
    sec_short_name                  varchar(128)   NULL,
    sec_name                        varchar(256)   NULL,
    sec_code                        varchar(12)    NULL,
    isin                            varchar(12)    NULL,
    class_code                      varchar(12)    NULL,
    operation                       varchar(16)    NOT NULL,
    trade_account                   varchar(12)    NULL,
    price                           decimal(19, 8) NULL,
    qty                             decimal(19, 8) NULL,
    trade_value                     decimal(19, 8) NULL,
    currency_code                   varchar(4)     NULL,
    settle_currency_code            varchar(4)     NULL,
    accrued_int                     decimal(19, 8) NULL,
    trader_code                     varchar(12)    NULL,
    firm_code                       varchar(12)    NULL,
    trader_firm_name                varchar(128)   NULL,
    client_code                     varchar(12)    NOT NULL,
    uid                             bigint         NULL,
    broker_ref                      varchar(128)   NULL,
    ts_commission                   decimal(19, 8) NULL,
    clearing_commission             decimal(19, 8) NULL,
    exchange_commission             decimal(19, 8) NULL,
    tech_center_commission          decimal(19, 8) NULL,
    broker_commission               decimal(19, 8) NULL,
    trade_kind                      varchar(128)   NULL,
    base_currency_code              varchar(4)     NULL,
    base_qty                        decimal(19, 8) NULL,
    quote_currency_code             varchar(4)     NULL,
    quote_qty                       decimal(19, 8) NULL,
    exchange_order_num              varchar(64)    NULL,
    exec_market                     varchar(64)    NULL,
    cross_rate                      decimal(19, 8) NULL,
    nominal_value                   decimal(19, 8) NULL,
    settle_value                    decimal(19, 8) NULL,
    ts_commission_currency_code     varchar(4)     NULL,
    broker_commission_currency_code varchar(4)     NULL,
    CONSTRAINT PK_quik_trades PRIMARY KEY CLUSTERED (
        client_code,
        operation,
        trade_num
    )
);

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Сделки по счетам клиента. Источник QUIK: окно «Сделки», выгрузка по ODBC. Обезличенные сделки сюда не входят.',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'trades';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Регистрационный номер сделки в торговой системе (QUIK: «Номер»). Часть ключа вместе с кодом клиента и направлением.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_num';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Идентификатор биржи (QUIK: «Код биржи»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'exchange_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата текущей торговой сессии (QUIK: «Дата торгов»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'session_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата регистрации сделки (QUIK: «Дата сделки»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата расчетов по сделке (QUIK: «Дата расчетов»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'settle_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Время регистрации сделки с точностью до секунды (QUIK: «Время»). Хранится как его отдал терминал.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_time';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Микросекунды времени регистрации сделки (QUIK: «Время (мкс)»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_time_mcs';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Номер заявки, на основании которой заключена сделка (QUIK: «Заявка»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'order_num';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Сокращенное наименование инструмента (QUIK: «Инструмент сокр.»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'sec_short_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Наименование инструмента, как его показывает терминал (QUIK: «Инструмент»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'sec_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код инструмента (QUIK: «Код инструмента»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'sec_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код ISIN.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'isin';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код класса (QUIK: «Код класса»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'class_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Направление сделки (QUIK: «Операция»): «Купля», «Продажа», «К/П», «П/К». Часть ключа. «К/П» и «П/К» не сводятся к купле или продаже.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'operation';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код торгового счета (QUIK: «Счет»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_account';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Цена сделки за единицу инструмента (QUIK: «Цена»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'price';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Количество инструментов (QUIK: «Кол-во»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'qty';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Объем сделки в денежном выражении (QUIK: «Объем»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_value';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта цены (QUIK: «Валюта»). Рубль приходит как SUR.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта расчетов (QUIK: «Валюта расчетов»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'settle_currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Накопленный купонный доход на количество в сделке (QUIK: «НКД»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'accrued_int';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Идентификатор трейдера (QUIK: «Трейдер»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trader_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Идентификатор фирмы, от имени которой совершена сделка (QUIK: «Дилер»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'firm_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Значение колонки «Орг-я трейдера». Руководство называет его идентификатором фирмы трейдера. В образце выгрузки это название. Хранится как пришло.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trader_firm_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код клиента (QUIK: «Код клиента»). Часть ключа.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'client_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код пользователя на сервере QUIK (QUIK: «UID»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'uid';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Комментарий трейдера (QUIK: «Комментарий»), обычно код клиента и номер поручения.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'broker_ref';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Комиссия торговой системы (QUIK: «Комиссия ТС»). На Мосбирже руководство задает ее как сумму клиринговой, ФБ и ТЦ комиссий. Сумма в базе не пересчитывается.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'ts_commission';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Клиринговая комиссия (QUIK: «Клиринговая комиссия», QLua: clearing_comission).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'clearing_commission';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Комиссия Фондовой биржи (QUIK: «ФБ комиссия»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'exchange_commission';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Комиссия Технического центра (QUIK: «ТЦ комиссия»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'tech_center_commission';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Комиссия брокера (QUIK: «Комиссия брокера»). В QLua поле резервное, точность до 2 знаков. В руководстве 2023 это сумма комиссии брокера и торговой системы в валюте фактического списания. Хранится сырое значение.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'broker_commission';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Вид сделки (QUIK: «Вид сделки»), например «Обычная». Это не признак маржинальной сделки из колонки «Тип сделки».', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'trade_kind';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Базовая валюта валютной пары (QUIK: «Базовая валюта»). Заполняется для класса «Валюта».', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'base_currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Количество базовой валюты (QUIK: «Кол-во базовой валюты»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'base_qty';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Котируемая валюта (QUIK: «Котируемая валюта»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'quote_currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Количество котируемой валюты (QUIK: «Кол-во котируемой валюты»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'quote_qty';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Биржевой номер заявки (QUIK: «Биржевой номер заявки»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'exchange_order_num';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Площадка исполнения (QUIK: «Площадка исполнения»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'exec_market';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Кросс-курс валюты цены к валюте расчетов (QUIK: «Кросс-курс»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'cross_rate';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Номинальный объем (QUIK: «Номинальный объем»). В руководстве это номинал инструмента, умноженный на размер лота и количество лотов. Хранится сырое значение.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'nominal_value';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Расчетная сумма (QUIK: «Расчетная сумма»). В просмотренном руководстве 2023 эта подпись не найдена. Формула не фиксируется, хранится сырое значение.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'settle_value';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта комиссии торговой системы (QUIK: «Валюта комиссии торговой системы»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'ts_commission_currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта комиссии брокера (QUIK: «Валюта комиссии брокера»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'trades', @level2type = N'COLUMN', @level2name = N'broker_commission_currency_code';

GRANT SELECT, INSERT, UPDATE, DELETE ON quik.trades TO quik_odbc_writer;
GRANT SELECT ON quik.trades TO quik_portfolio_reader;

-- +goose Down

DROP TABLE quik.trades;
