package quik

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JohannesJHN/iso4217"
)

var (
	ErrNotExistingCurrency = errors.New("валюта не существует в ISO 4217")
)

type currencyIso struct {
	alpha      CurrencyCode
	numeric    int16
	latName    string
	minorUnits int32
}

func (c currencyIso) LatName() string {
	return c.latName
}
func (c currencyIso) MinorUnits() int32 {
	return c.minorUnits
}
func (c currencyIso) Numeric() int16 {
	return c.numeric
}

func (c currencyIso) Alpha() CurrencyCode {
	return c.alpha
}

type Currency struct {
	charCode string
	name     *string
	currencyIso
}

func (c Currency) CharCode() string {
	return c.charCode
}

func (c Currency) Name() *string {
	return c.name
}

func CurrencyFromQuik(charCode string, name *string) (Currency, error) {

	iso, err := resolveCurrencyIso(charCode)
	if err != nil {
		return Currency{}, err
	}
	ccy := Currency{
		currencyIso: iso,
		name:        name,
		charCode:    charCode,
	}
	return ccy, nil
}

func resolveCurrencyIso(rawCode string) (currencyIso, error) {
	upper := strings.ToUpper(strings.TrimSpace(rawCode))

	alpha, ok := alphaByQuikCode[upper]
	if !ok {
		alpha = CurrencyCode(upper)
		if err := alpha.validateFormat(); err != nil {
			return currencyIso{}, err
		}

	}
	iso, ok := iso4217.LookupByAlpha3(alpha.String())
	if !ok {
		return currencyIso{}, fmt.Errorf("%w: код=%s", ErrNotExistingCurrency, alpha)
	}
	var miu int32 = int32(iso.MinorUnits)
	if miu < 0 {
		miu = 0
	}

	return currencyIso{
		alpha:      alpha,
		numeric:    int16(iso.Numeric),
		latName:    iso.Name,
		minorUnits: miu,
	}, nil
}
