package service

import (
	"fmt"
	"strings"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func normalizeQuikCcy(code string) string {
	upper := strings.ToUpper(strings.TrimSpace(code))
	if upper == "SUR" || upper == "RUR" {
		return "RUB"
	}
	return upper
}

func validateCurrencyCode(code string) error {
	normalized := normalizeQuikCcy(code)
	if err := currencies.CheckCurrencyCode(normalized); err != nil {
		return fmt.Errorf("%w: %s", md.ErrBusinessValidation, err.Error())
	}
	return nil
}
