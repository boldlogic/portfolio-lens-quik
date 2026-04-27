package currencies

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JohannesJHN/iso4217"
)

var (
	ErrWrongISOCharCode    = errors.New("символьный код валюты по ISO 4217 состоит из 3 латинских букв")
	ErrNotExistingCurrency = errors.New("валюта не существует в ISO 4217")
)

func CheckCurrencyCode(charCode string) error {
	code := strings.ToUpper(charCode)
	if len(code) != 3 {
		return ErrWrongISOCharCode
	}

	for _, r := range code {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return ErrWrongISOCharCode
		}
	}
	if _, ok := iso4217.LookupByAlpha3(code); !ok {
		return fmt.Errorf("%w: код=%s", ErrNotExistingCurrency, code)
	}
	return nil
}
