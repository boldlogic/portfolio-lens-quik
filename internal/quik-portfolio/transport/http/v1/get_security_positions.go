package v1

import (
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type positionDTO struct {
	ClientCode                      string          `json:"clientCode"`
	FirmCode                        string          `json:"firmCode"`
	FirmName                        string          `json:"firmName,omitempty"`
	Ticker                          string          `json:"ticker,omitempty"`
	Name                            string          `json:"name,omitempty"`
	Amount                          decimal.Decimal `json:"amount"`
	MarketValueInInstrumentCurrency decimal.Decimal `json:"marketValueInInstrumentCurrency"`
	InstrumentCurrencyCode          string          `json:"instrumentCurrencyCode,omitempty"`
	MarketValueInTargetCurrency     decimal.Decimal `json:"marketValueInTargetCurrency"`
}

func positionToDTO(p quik.Position) positionDTO {
	out := positionDTO{
		ClientCode:                      p.ClientCode,
		FirmCode:                        p.FirmCode,
		FirmName:                        p.FirmName,
		Ticker:                          p.Ticker,
		Name:                            p.Name,
		Amount:                          p.Amount,
		MarketValueInTargetCurrency:     p.MarketValueInTargetCurrency,
		MarketValueInInstrumentCurrency: p.MarketValueInInstrCurrency,
		InstrumentCurrencyCode:          p.InstrumentCurrencyCode,
	}

	return out
}
