package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type portfolioEntryDTO struct {
	LimitType      string          `json:"limitType"`
	LoadDate       string          `json:"loadDate"`
	SourceDate     string          `json:"sourceDate"`
	ClientCode     string          `json:"clientCode"`
	FirmCode       string          `json:"firmCode,omitempty"`
	FirmName       string          `json:"firmName"`
	Instrument     string          `json:"instrument"`
	ISIN           string          `json:"isin,omitempty"`
	ShortName      string          `json:"shortName,omitempty"`
	QuoteDate      string          `json:"quoteDate,omitempty"`
	QTY            decimal.Decimal `json:"qty"`
	MvCurrency     string          `json:"mvCurrency,omitempty"`
	MvInCcy        decimal.Decimal `json:"mvInCcy"`
	MvPrice        decimal.Decimal `json:"mvPrice"`
	MvAccrued      decimal.Decimal `json:"mvAccrued"`
	MvTotal        decimal.Decimal `json:"mvTotal"`
	TargetCurrency string          `json:"targetCurrency,omitempty"`
}

func (h *Handler) getPortfolio(r *http.Request) (any, string, error) {
	ctx := r.Context()

	targetCcy := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get(currencyQuery)))
	if targetCcy == "" {
		targetCcy = "RUB"
	}

	entries, err := h.service.GetPortfolio(ctx, targetCcy)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}

	return portfolioEntriesToDTO(entries), "", nil
}

func portfolioEntriesToDTO(entries []quik.PortfolioEntry) []portfolioEntryDTO {
	result := make([]portfolioEntryDTO, 0, len(entries))
	for _, e := range entries {
		dto := portfolioEntryDTO{
			LimitType:  string(e.LimitType),
			LoadDate:   e.LoadDate.Format(dates.ISODateFormat),
			SourceDate: e.SourceDate.Format(dates.ISODateFormat),
			ClientCode: e.ClientCode,
			FirmCode:   e.FirmCode,
			FirmName:   e.FirmName,
			Instrument: e.Instrument,
			QTY:        e.Balance,
			MvInCcy:    e.MvInCcy,
			MvPrice:    e.MvPrice,
			MvAccrued:  e.MvAccrued,
			MvTotal:    e.MvTotal,
		}
		if e.ISIN != nil {
			dto.ISIN = *e.ISIN
		}
		if e.ShortName != nil {
			dto.ShortName = *e.ShortName
		}
		dto.MvCurrency = e.MvCurrency
		dto.TargetCurrency = e.TargetCurrency
		if e.QuoteDate != nil {
			dto.QuoteDate = e.QuoteDate.Format(dates.ISODateFormat)
		}
		result = append(result, dto)
	}
	return result
}
