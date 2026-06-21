package quik

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

type LimitType string

const (
	LimitTypeSecurities    LimitType = "security"     // ценные бумаги (биржевые)
	LimitTypeSecuritiesOtc LimitType = "security_otc" // ценные бумаги OTC
	LimitTypeMoney         LimitType = "money"        // денежные средства
)

const (
	maxClientCodeLen   = 12
	minClientCodeLen   = 1
	maxFirmCodeLen     = 12
	minFirmCodeLen     = 1
	maxPositionCodeLen = 4
	minPositionCodeLen = 1
	maxTradeAccountLen = 12
	minTradeAccountLen = 1
	maxTickerLen       = 12
	minTickerLen       = 1
)

var ErrWrongLimitType = errors.New("тип лимита должен быть money, security или security_otc")

func (lt LimitType) Validate() error {
	switch lt {
	case LimitTypeSecurities, LimitTypeSecuritiesOtc, LimitTypeMoney:
		return nil
	default:
		return ErrWrongLimitType
	}
}

type Limit struct {
	limitType               LimitType
	clientCode              string     // код клиента
	ticker                  string     // валюта для ML
	settleCode              SettleCode // Срок расчетов
	positionCode            *string    // Код позиции для ML
	tradeAccount            *string
	firmCode                string //
	balance                 decimal.Decimal
	acquisitionCurrencyCode *string
	isin                    *string
}

func NewLimit(limitType string,
	clientCode string,
	ticker string,
	positionCode *string,
	settleCode string,
	tradeAccount *string,
	firmCode string,
	balance decimal.Decimal,
	acquisitionCurrencyCode *string,
	ISIN *string,
) (Limit, error) {
	if settleCode == "" {
		settleCode = SettleCodeTx
	}
	out := Limit{
		limitType:  LimitType(strings.TrimSpace(limitType)),
		clientCode: strings.TrimSpace(clientCode),
		ticker:     strings.TrimSpace(ticker),
		firmCode:   strings.TrimSpace(firmCode),
		balance:    balance,
		settleCode: SettleCode(settleCode),
	}
	err := out.validateCommon()
	if err != nil {
		return Limit{}, err
	}

	switch out.limitType {
	case LimitTypeMoney:
		out.positionCode = positionCode

		err = out.validateMoney()
		if err != nil {
			return Limit{}, err
		}

	case LimitTypeSecuritiesOtc:
		if tradeAccount == nil || strings.TrimSpace(*tradeAccount) == "" {
			tradeAccount = new("OTC")
		}

		fallthrough
	case LimitTypeSecurities:
		out.tradeAccount = tradeAccount

		if acquisitionCurrencyCode != nil && strings.TrimSpace(*acquisitionCurrencyCode) != "" {
			out.acquisitionCurrencyCode = acquisitionCurrencyCode
		}
		if ISIN != nil && strings.TrimSpace(*ISIN) != "" {
			out.isin = ISIN
		}

		err = out.validateSecurity()

		if err != nil {
			return Limit{}, err
		}

	}

	return out, nil
}

func (l Limit) Type() LimitType {
	return l.limitType
}

func (l Limit) ClientCode() string {
	return l.clientCode
}

func (l Limit) Ticker() string {
	return l.ticker
}

func (l Limit) SettleCode() SettleCode {
	return l.settleCode
}

func (l Limit) PositionCode() *string {
	return l.positionCode
}

func (l Limit) TradeAccount() *string {
	return l.tradeAccount
}

func (l Limit) FirmCode() string {
	return l.firmCode
}

func (l Limit) Balance() decimal.Decimal {
	return l.balance
}

func (l Limit) AcquisitionCurrencyCode() *string {
	return l.acquisitionCurrencyCode
}

func (l Limit) ISIN() *string {
	return l.isin
}

func (l Limit) String() string {
	return fmt.Sprintf(
		"Limit{Type:%q ClientCode:%q Ticker:%q SettleCode:%q PositionCode:%s TradeAccount:%s FirmCode:%q Balance:%s AcquisitionCurrencyCode:%s ISIN:%s}",
		l.Type(),
		l.ClientCode(),
		l.Ticker(),
		l.SettleCode(),
		limitOptionalString(l.PositionCode()),
		limitOptionalString(l.TradeAccount()),
		l.FirmCode(),
		l.Balance().String(),
		limitOptionalString(l.AcquisitionCurrencyCode()),
		limitOptionalString(l.ISIN()),
	)
}

func limitOptionalString(s *string) string {
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf("%q", *s)
}

func ParseClientCode(raw string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(raw))
	if err := validateClientCodeLength(code); err != nil {
		return "", err
	}
	return code, nil
}

func validateClientCodeLength(code string) error {
	cl := utf8.RuneCountInString(code)
	if cl < minClientCodeLen || cl > maxClientCodeLen {
		return fmt.Errorf("код клиента должен иметь длину от %d до %d символов", minClientCodeLen, maxClientCodeLen)
	}
	return nil
}

func (l Limit) validateCommon() error {
	err := l.limitType.Validate()
	if err != nil {
		return err
	}

	if err := validateClientCodeLength(l.clientCode); err != nil {
		return err
	}
	fc := utf8.RuneCountInString(l.firmCode)
	if fc < minFirmCodeLen || fc > maxFirmCodeLen {
		return fmt.Errorf("код фирмы должен иметь длину от %d до %d символов", minFirmCodeLen, maxFirmCodeLen)
	}

	return l.settleCode.Validate()
}

func (l Limit) validateMoney() error {

	_, err := ParseCurrencyCode(l.ticker)
	if err != nil {
		return err
	}
	if l.positionCode == nil {
		return fmt.Errorf("не указан код позиции")
	}
	pc := utf8.RuneCountInString(*l.positionCode)
	if pc < minPositionCodeLen || pc > maxPositionCodeLen {
		return fmt.Errorf("код позиции должен иметь длину от %d до %d символов", minPositionCodeLen, maxPositionCodeLen)

	}

	return nil

}

func (l Limit) validateSecurity() error {

	pc := utf8.RuneCountInString(l.ticker)
	if pc < minTickerLen || pc > maxTickerLen {
		return fmt.Errorf("код бумаги должен иметь длину от %d до %d символов", minTickerLen, maxTickerLen)

	}

	if l.tradeAccount == nil {
		return fmt.Errorf("не указан торговый счет")
	}
	ta := utf8.RuneCountInString(*l.tradeAccount)
	if ta < minTradeAccountLen || ta > maxTradeAccountLen {
		return fmt.Errorf("код позиции должен иметь длину от %d до %d символов", minTradeAccountLen, maxTradeAccountLen)
	}

	if l.isin != nil {
		i := utf8.RuneCountInString(*l.isin)
		if i < minTickerLen || i > maxTickerLen {
			return fmt.Errorf("isin бумаги должен иметь длину от %d до %d символов", minTickerLen, maxTickerLen)
		}

	}
	if l.acquisitionCurrencyCode != nil {
		if utf8.RuneCountInString(*l.acquisitionCurrencyCode) > 4 {
			return fmt.Errorf("валюта приобретения не может иметь длину более 4 символов")
		}
	}

	return nil

}

func appendLimitKeyPart(b []byte, part string) []byte {
	if len(b) > 0 {
		b = append(b, 0)
	}
	return append(b, part...)
}

func (l Limit) keyBytes() ([]byte, bool) {
	var b []byte
	switch l.limitType {
	case LimitTypeMoney:
		if l.positionCode == nil {
			return nil, false
		}
		b = appendLimitKeyPart(b, string(l.limitType))
		b = appendLimitKeyPart(b, l.clientCode)
		b = appendLimitKeyPart(b, l.ticker)
		b = appendLimitKeyPart(b, string(l.settleCode))
		b = appendLimitKeyPart(b, *l.positionCode)
		b = appendLimitKeyPart(b, l.firmCode)
	case LimitTypeSecurities, LimitTypeSecuritiesOtc:
		if l.tradeAccount == nil {
			return nil, false
		}
		b = appendLimitKeyPart(b, string(l.limitType))
		b = appendLimitKeyPart(b, l.clientCode)
		b = appendLimitKeyPart(b, l.ticker)
		b = appendLimitKeyPart(b, string(l.settleCode))
		b = appendLimitKeyPart(b, *l.tradeAccount)
		b = appendLimitKeyPart(b, l.firmCode)
	default:
		return nil, false
	}
	return b, true
}

func (l Limit) KeyHash() [32]byte {
	b, ok := l.keyBytes()
	if !ok {
		return [32]byte{}
	}
	return sha256.Sum256(b)
}
