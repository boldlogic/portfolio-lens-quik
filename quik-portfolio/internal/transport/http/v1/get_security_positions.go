package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
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
	MarketValueInTargetCurrency     decimal.Decimal `json:"marketValueInTargetCurrency"`
	InstrumentCurrencyCode          string          `json:"instrumentCurrencyCode,omitempty"`
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
	// switch p.LimitType {
	// case quik.LimitTypeMoney:
	// 	return out
	// default:
	// 	out.InstrumentCurrencyCode = p.InstrumentCurrencyCode

	// }
	return out
}

func (h *Handler) getSecurityPositions(r *http.Request) (any, string, error) {
	ctx := r.Context()
	params := parsePortfolioQueryParams(r)

	positions, total, portfolioCCY, err := h.service.GetSecurityPositions(ctx, params.TargetCurrency, params.ClientCodes)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return positionsToPortfolioDTO(positions, total, portfolioCCY), "", nil
}
