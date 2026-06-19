package quik

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/shopspring/decimal"
)

type LimitType string

const (
	LimitTypeSecurities    LimitType = "security"     // ценные бумаги (биржевые)
	LimitTypeSecuritiesOtc LimitType = "security_otc" // ценные бумаги OTC
	LimitTypeMoney         LimitType = "money"        // денежные средства
)

var ErrWrongLimitType = errors.New("LimitType должен быть security, security_otc или money")

func (lt LimitType) Validate() error {
	switch lt {
	case LimitTypeSecurities, LimitTypeSecuritiesOtc, LimitTypeMoney:
		return nil
	default:
		return ErrWrongLimitType
	}
}

const (
	MaxClientCodeLen = 12
	MinClientCodeLen = 1
)

type Limit struct {
	Type                    LimitType
	LoadDate                time.Time // Дата лимита
	SourceDate              time.Time // Дата загрузки
	ClientCode              string    // код клиента
	CurrencyCode            *string   // валюта для ML
	SecCode                 *string
	PositionCode            *string    // Код позиции для ML
	SettleCode              SettleCode // Срок расчетов
	TradeAccount            *string
	FirmCode                string //
	FirmName                string
	Balance                 decimal.Decimal
	AcquisitionCurrencyCode *string
	ISIN                    *string
	ShortName               *string
}

func NewLimit(limitType string,
	clientCode string, // код клиента
	ticker *string,
	positionCode *string, // Код позиции для ML
	settleCode string, // Срок расчетов
	tradeAccount *string,
	firmCode string, //
	firmName *string,
	balance decimal.Decimal,
	acquisitionCurrencyCode *string,
	ISIN *string,
	shortName *string,
) (*Limit, error) {
	tp := LimitType(limitType)
	err := tp.Validate()
	if err != nil {
		return nil, err
	}
	out := Limit{
		Type:       tp,
		FirmCode:   firmCode,
		ClientCode: clientCode,
	}

	if strings.TrimSpace(settleCode) == "" {
		settleCode = SettleCodeTx
	}
	sc := SettleCode(settleCode)
	if err := sc.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
	}
	out.SettleCode = sc

	switch tp {
	case LimitTypeMoney:
		if ticker == nil || strings.TrimSpace(*ticker) == "" {
			return nil, fmt.Errorf("%w: не указана валюта", models.ErrBusinessValidation)
		}
		normCcy, err := ParseCurrencyCode(*ticker)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
		}
		ccy := normCcy.String()
		out.CurrencyCode = &ccy
		if positionCode == nil || strings.TrimSpace(*positionCode) == "" {
			return nil, fmt.Errorf("%w: не указан код позиции", models.ErrBusinessValidation)
		}
		out.PositionCode = positionCode
	case LimitTypeSecuritiesOtc:
		if tradeAccount == nil || strings.TrimSpace(*tradeAccount) == "" {
			otc := "OTC"
			tradeAccount = &otc
		}
		out.TradeAccount = tradeAccount

		fallthrough
	case LimitTypeSecurities:
		if ticker == nil || strings.TrimSpace(*ticker) == "" {
			return nil, fmt.Errorf("%w: не указан код бумаги", models.ErrBusinessValidation)
		}
		out.SecCode = ticker
		if tradeAccount == nil || strings.TrimSpace(*tradeAccount) == "" {
			return nil, fmt.Errorf("%w: не указан торговый счет", models.ErrBusinessValidation)
		}
		out.TradeAccount = tradeAccount
		out.ISIN = ISIN
		out.ShortName = shortName
	}

	out.Balance = balance
	return &out, nil
}

type MoneyLimit struct {
	LoadDate     time.Time
	SourceDate   time.Time
	ClientCode   string
	CurrencyCode string
	PositionCode string
	SettleCode   SettleCode
	FirmCode     string
	FirmName     string
	Balance      decimal.Decimal
}

type SecurityLimit struct {
	LoadDate                time.Time
	SourceDate              time.Time
	ClientCode              string
	SecCode                 string
	TradeAccount            string
	SettleCode              SettleCode
	FirmCode                string
	FirmName                string
	Balance                 decimal.Decimal
	AcquisitionCurrencyCode string
	ISIN                    string
	ShortName               string
}

type Position struct {
	LimitType                   LimitType
	LoadDate                    time.Time       //дата исходного лимита
	SourceDate                  time.Time       // если не равен LoadDate то означает дату, с которой перенесен лимит
	ClientCode                  string          //код клиента
	FirmCode                    string          //код фирмы
	FirmName                    string          //название фирмы
	Ticker                      string          //код инструмента
	Name                        string          //название инструмента
	Amount                      decimal.Decimal //текущее фактическое количество
	UnitPrice                   decimal.Decimal //цена позиции
	AccruedInterest             decimal.Decimal //НКД в валюте инструмента
	MarketValueInInstrCurrency  decimal.Decimal //оценка позиции в валюте инструмента
	MarketValueInTargetCurrency decimal.Decimal // оценка позиции в валюте запроса
	InstrumentCurrencyCode      string          //валюта инструмента
}
