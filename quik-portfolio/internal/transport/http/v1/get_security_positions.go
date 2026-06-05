package v1

import (
	"errors"
	"net/http"
	"strings"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type position struct {
	ClientCode   string          `json:"clientCode"`
	FirmCode     string          `json:"firmCode"`
	FirmName     string          `json:"firmName,omitempty"`
	Ticker       string          `json:"ticker,omitempty"`
	Name         string          `json:"name,omitempty"`
	Amount       decimal.Decimal `json:"amount"`
	MVInstr      decimal.Decimal `json:"marketValueInstrument"`
	MV           decimal.Decimal `json:"marketValue"`
	CurrencyCode string          `json:"currency,omitempty"`
}

func toDTO(p quik.Position) position {
	out := position{
		ClientCode: p.ClientCode,
		FirmCode:   p.FirmCode,
		FirmName:   p.FirmName,
		Ticker:     p.Ticker,
		Name:       p.Name,
		Amount:     p.Balance,
		MV:         p.MVTotal,
		MVInstr:    p.MVInstr,
	}
	switch p.LimitType {
	case quik.LimitTypeMoney:
		return out
	default:
		out.CurrencyCode = p.CurrencyCode
		return out
	}

}

func (h *Handler) getSecurityPositions(r *http.Request) (any, string, error) {
	ctx := r.Context()
	var currency *string
	targetCcy := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get(currencyQuery)))
	if targetCcy != "" {
		currency = &targetCcy
	}

	positions, total, portfolioCCY, err := h.service.GetSecurityPositions(ctx, currency)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return portfolioToDTO(positions, total, portfolioCCY), "", nil
}
