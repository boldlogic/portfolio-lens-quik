-- +goose Up

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Схема объектов, заполняемых из терминала QUIK по ODBC.',
    @level0type = N'SCHEMA', @level0name = N'quik';

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Лимиты по денежным средствам. Источник QUIK: таблица MoneyLimits, окно «Позиции по деньгам».',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'money_limits';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата среза позиции.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'load_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код клиента (QUIK: ClientCode).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'client_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код валюты (QUIK: CurrCode).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код позиции (QUIK: Tag).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'position_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Срок расчетов (QUIK: LimitKind).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'settle_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код участника торгов (QUIK: FirmId).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'firm_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Наименование участника торгов (QUIK: FirmName).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'firm_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Текущий остаток денежных средств (QUIK: CurrentBal).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'balance';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата загрузки лимита в таблицу.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'source_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Служебная метка строки для контроля изменений.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'money_limits', @level2type = N'COLUMN', @level2name = N'ts';

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Лимиты по ценным бумагам (биржевые). Источник QUIK: таблица DepoLimits, окно «Позиции по инструментам».',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'security_limits';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата среза позиции.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'load_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код клиента (QUIK: ClientCode).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'client_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код финансового инструмента (QUIK: SecCode).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'sec_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Счет депо / торговый счет (QUIK: TrdAcc).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'trade_account';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Срок расчетов / тип лимита (QUIK: LimitKind).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'settle_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код участника торгов (QUIK: FirmId).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'firm_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Наименование участника торгов (QUIK: FirmName).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'firm_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Текущий остаток инструментов в лотах или штуках (QUIK: CurrentBal).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'balance';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта цены приобретения. Для облигаций QUIK может передавать «%» (QUIK: WAPositionPriceCurrency).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'acquisition_currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Международный идентификатор ценной бумаги (ISIN).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'isin';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Полное наименование инструмента (QUIK: SecName).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'sec_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата загрузки лимита в таблицу.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'source_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Служебная метка строки для контроля изменений.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits', @level2type = N'COLUMN', @level2name = N'ts';

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Лимиты по ценным бумагам OTC. QUIK не выгружает; загрузка через API сервиса writer.',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'security_limits_otc';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата среза позиции.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'load_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код клиента.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'client_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код финансового инструмента.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'sec_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Торговый счет; для OTC по умолчанию OTC.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'trade_account';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Срок расчетов / тип лимита.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'settle_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код участника торгов.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'firm_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Наименование участника торгов.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'firm_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Текущий остаток инструментов.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'balance';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта цены приобретения.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'acquisition_currency_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Международный идентификатор ценной бумаги (ISIN).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'isin';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Наименование инструмента.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'sec_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата загрузки лимита в таблицу.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'source_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Служебная метка строки для контроля изменений.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'security_limits_otc', @level2type = N'COLUMN', @level2name = N'ts';

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Текущие котировки и параметры инструментов. Источник QUIK: таблица Params, окно «Текущие торги».',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'current_quotes';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата котировки.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'quote_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Составной ключ инструмента и класса (QUIK: «Инструмент + Класс»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'instrument_class';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код инструмента (QUIK: SecCode).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'sec_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Регистрационный номер инструмента (QUIK: Securities.REGNUMBER).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'registration_number';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Полное наименование инструмента (QUIK: Securities.SHORTNAME / «Инструмент»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'full_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Краткое наименование инструмента (QUIK: «Инструмент сокр.»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'short_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код класса / борда (QUIK: ClassCode).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'class_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Наименование класса (QUIK: ClassName / «Класс»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'class_name';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Тип инструмента (QUIK: «Тип инстр-та» / Securities.SECTYPE).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'instrument_type';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Подтип инструмента (QUIK: «Подтип инстр-та»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'instrument_subtype';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Международный идентификатор ценной бумаги (ISIN).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'isin';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Номинал инструмента (QUIK: FaceValue).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'face_value';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта инструмента (QUIK: «Валюта»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'currency';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Базовая валюта (QUIK: BaseCurrency).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'base_currency';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Валюта котировки (QUIK: «Котир.валюта»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'quote_currency';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Сопряженная валюта (QUIK: «Сопр.валюта»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'counter_currency';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Дата погашения (QUIK: MaturityDate / «Погашение»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'maturity_date';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Длительность купона в днях (QUIK: couponperiod / «Длит. купона»).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'coupon_duration';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Цена последней сделки (QUIK: last).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'last_price';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Цена закрытия периода (QUIK: closeprice).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'close_price';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Средневзвешенная цена (QUIK: waprice).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'waprice';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Накопленный купонный доход (QUIK: accruedint).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'accrued_int';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Статус торгов (QUIK: tradingstatus).', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'trading_status';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Идентификатор инструмента в справочнике системы; заполняется сервисом instruments.', @level0type = N'SCHEMA', @level0name = N'quik', @level1type = N'TABLE', @level1name = N'current_quotes', @level2type = N'COLUMN', @level2name = N'instrument_id';
