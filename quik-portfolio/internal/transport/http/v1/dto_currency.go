package v1

import (
	"strings"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
)

func currencyCodeForDTO(raw string) string {
	trimmed := strings.ToUpper(strings.TrimSpace(raw))
	if trimmed == "" {
		return ""
	}
	norm, err := currencies.ParseCurrencyCode(trimmed)
	if err == nil {
		return norm.String()
	}
	switch trimmed {
	case "SUR", "RUR":
		return "RUB"
	default:
		return trimmed
	}
}
