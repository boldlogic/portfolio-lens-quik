package quik

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/JohannesJHN/iso4217"
)

var alphaByQuikCode = map[string]CurrencyCode{
	"GLD":  "XAU",
	"SUR":  "RUB",
	"RUR":  "RUB",
	"USDX": "USD",
	"SLV":  "XAG",
	"PLT":  "XPT",
	"PLD":  "XPD",
}

type CurrencyCode string

func (c CurrencyCode) String() string {
	return string(c)
}

var (
	ErrWrongCurrencyCode = errors.New("символьный код валюты по ISO 4217 состоит из 3 латинских букв")
)

func ParseCurrencyCode(rawCode string) (CurrencyCode, error) {
	upper := strings.ToUpper(strings.TrimSpace(rawCode))
	code, ok := alphaByQuikCode[upper]
	if ok {
		return code, nil
	}

	code = CurrencyCode(upper)
	if err := code.Validate(); err != nil {
		return "", err
	}

	return code, nil
}

func (c CurrencyCode) Validate() error {
	if utf8.RuneCountInString(c.String()) != 3 {
		return ErrWrongCurrencyCode
	}
	for _, r := range c {
		if r < 'A' || r > 'Z' {
			return ErrWrongCurrencyCode
		}
	}
	_, ok := iso4217.LookupByAlpha3(c.String())
	if !ok {
		return ErrNotExistingCurrency
	}
	return nil
}

///
