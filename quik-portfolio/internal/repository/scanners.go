package repository

var (
	ErrScan     = errors.New("ошибка чтения строки")
)


//+
type moneyPosition struct {
	LoadDate     time.Time
	SourceDate   time.Time
	RateDate     sql.NullTime
	ClientCode   string
	FirmCode     string
	FirmName     string
	CurrencyCode string
	CurrencyName sql.NullString
	Balance      decimal.Decimal
	MV           decimal.Decimal
}
//+
func scanMoneyToPosition(row *sql.Rows) (moneyPosition, error) {
	var out moneyPosition
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.RateDate,
		&out.ClientCode,
		&out.FirmCode,
		&out.FirmName,
		&out.CurrencyCode,
		&out.CurrencyName,
		&out.Balance,
		&out.MV,
	)
	if err != nil {
		return moneyPosition{}, fmt.Errorf("%w: %w",ErrScan,err)
	}

	return out, nil
}

func scanMoneyLimit(row *sql.Rows) (quik.MoneyLimit, error) {
	out := quik.MoneyLimit{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.Currency,
		&out.PositionCode,
		&out.SettleCode,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
	)
	if err != nil {
		return quik.MoneyLimit{}, fmt.Errorf("%w: %w",ErrScan,err)
	}
	return out, nil
}

func scanSecurityLimit(row *sql.Rows) (quik.SecurityLimit, error) {
	out := quik.SecurityLimit{}
	var shortName sql.NullString
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.Ticker,
		&out.SettleCode,
		&out.TradeAccount,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
		&out.AcquisitionCcy,
		&out.ISIN,
		&shortName,
	)
	if err != nil {
		return quik.SecurityLimit{}, fmt.Errorf("%w: %w",ErrScan,err)
	}
	if shortName.Valid {
		out.ShortName = shortName.String
	}
	return out, nil
}

func scanSecurityToPosition(row *sql.Rows) (securityPosition, error) {
	var out securityPosition
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.QuoteDate,
		&out.RateDate,
		&out.ClientCode,
		&out.FirmCode,
		&out.FirmName,
		&out.SecCode,
		&out.SecName,
		&out.Balance,
		&out.Price,
		&out.AccruedInt,
		&out.MVInstr,
		&out.CurrencyCode,
		&out.MV,
	)
	if err != nil {
		return securityPosition{}, fmt.Errorf("%w: %w",ErrScan,err)
	}

	return out, nil
}