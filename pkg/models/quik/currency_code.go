package quik

import (
	"errors"
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
	iso, err := resolveCurrencyIso(rawCode)
	if err != nil {
		return "", err
	}

	return iso.alpha, nil
}

func (c CurrencyCode) Validate() error {
	if err := c.validateFormat(); err != nil {
		return err
	}

	_, ok := iso4217.LookupByAlpha3(c.String())
	if !ok {
		return ErrNotExistingCurrency
	}
	return nil
}

func (c CurrencyCode) validateFormat() error {
	if utf8.RuneCountInString(c.String()) != 3 {
		return ErrWrongCurrencyCode
	}
	for _, r := range c {
		if r < 'A' || r > 'Z' {
			return ErrWrongCurrencyCode
		}
	}
	return nil
}
