package currencies

import (
	"errors"
	"testing"
)

func Test_CheckCurrencyCode_CharCodeOK(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"3 буквы A-Z", "USD"},
		{"3 буквы a-z", "eur"},
		{"смешанный регистр", "Usd"},
		{"RUB", "RUB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := CheckCurrencyCode(tt.code)
			if err != nil {
				t.Errorf("Test_CheckCurrencyCode_CharCodeOK(%q) = %v, ОР %v", tt.code, err, "nil")
			}

		})
	}
}
func Test_CheckCurrencyCode_WrongISOCharCode(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"2 буквы", "US"},
		{"4 буквы", "USDD"},
		{"цифры", "US1"},
		{"кириллица", "РУБ"},
		{"пусто", ""},
		{"пробел", "  "},
		{"спецсимвол", "U$D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := CheckCurrencyCode(tt.code)
			if !errors.Is(err, ErrWrongISOCharCode) {
				t.Errorf("Test_CheckCurrencyCode_WrongISOCharCode(%q) = %v, ОР %v", tt.code, err, ErrWrongISOCharCode)
			}

		})
	}
}

func Test_CheckCurrencyCode_ErrNotExistingCurrency(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"AAA", "AAA"},
		{"aaa", "aaa"},
		{"GLD", "GLD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := CheckCurrencyCode(tt.code)
			//ne := *ErrNotExistingCurrency
			if !errors.Is(err, ErrNotExistingCurrency) {
				t.Errorf("Test_CheckCurrencyCode_WrongISOCharCode(%q) = %v, ОР %v", tt.code, err, ErrNotExistingCurrency)
			}

		})
	}
}
