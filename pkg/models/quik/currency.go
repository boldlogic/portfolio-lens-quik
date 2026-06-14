package quik

import (
	"errors"

	"github.com/JohannesJHN/iso4217"
)

type currencyAlphaCode string

func (c currencyAlphaCode) String() string {
	return string(c)
}

var (
	ErrNotExistingCurrency = errors.New("валюта не существует в ISO 4217")
)

type currencyIso struct {
	alpha      CurrencyCode
	numeric    int16
	latName    string
	minorUnits int32
}

type Currency struct {
	charCode string
	name     *string
	*currencyIso
}

func (c Currency) CharCode() string {
	return c.charCode
}

func (c Currency) Name() *string {
	return c.name
}

func (c Currency) LatName() string {
	return c.latName
}
func (c Currency) MinorUnits() int32 {
	return c.minorUnits
}
func (c Currency) Numeric() int16 {
	return c.numeric
}

func (c Currency) Alpha() CurrencyCode {
	return c.alpha
}

func CurrencyFromQuik(charCode string, name *string) (Currency, error) {
	alpha, err := ParseCurrencyCode(charCode)
	if err != nil {
		return Currency{}, err
	}
	iso, _ := iso4217.LookupByAlpha3(alpha.String())

	var miu int32 = int32(iso.MinorUnits)
	if miu < 0 {
		miu = 0
	}

	ccyIso := currencyIso{
		alpha:      alpha,
		numeric:    int16(iso.Numeric),
		latName:    iso.Name,
		minorUnits: int32(iso.MinorUnits),
	}
	ccy := Currency{
		currencyIso: &ccyIso,
		name:        name,
		charCode:    charCode,
	}
	return ccy, nil
}

// func NewCurrency(charCode string, name *string) (*Currency, error) {

// 	alpha, err := newCurrencyAlphaCode(charCode)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ccy, ok := iso4217.LookupByAlpha3(alpha.String())
// 	if !ok {
// 		return nil, fmt.Errorf("%w: код=%s", ErrNotExistingCurrency, alpha)
// 	}
// 	var miu int32 = int32(ccy.MinorUnits)
// 	if miu < 0 {
// 		miu = 0
// 	}

// 	cur := Currency{
// 		Alpha:      alpha,
// 		Numeric:    int16(ccy.Numeric),
// 		CharCode:   charCode,
// 		LatName:    ccy.Name,
// 		MinorUnits: miu,
// 	}

// 	if name != nil {
// 		cur.Name = name
// 	}

// 	return &cur, nil
// }
