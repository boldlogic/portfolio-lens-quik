package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (h *Handler) GetPortfolio(r *http.Request) (any, string, error) {
	ctx := r.Context()

	targetCcy := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("targetCcy")))
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
			QTY:        e.Balance.InexactFloat64(),
			MvInCcy:    e.MvInCcy.InexactFloat64(),
			MvPrice:    e.MvPrice.InexactFloat64(),
			MvAccrued:  e.MvAccrued.InexactFloat64(),
			MvTotal:    e.MvTotal.InexactFloat64(),
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

type portfolioEntryDTO struct {
	LimitType      string  `json:"limitType"`
	LoadDate       string  `json:"loadDate"`
	SourceDate     string  `json:"sourceDate"`
	ClientCode     string  `json:"clientCode"`
	FirmCode       string  `json:"firmCode,omitempty"`
	FirmName       string  `json:"firmName"`
	Instrument     string  `json:"instrument"`
	ISIN           string  `json:"isin,omitempty"`
	ShortName      string  `json:"shortName,omitempty"`
	QuoteDate      string  `json:"quoteDate,omitempty"`
	QTY            float64 `json:"qty"`
	MvCurrency     string  `json:"mvCurrency,omitempty"`
	MvInCcy        float64 `json:"mvInCcy"`
	MvPrice        float64 `json:"mvPrice"`
	MvAccrued      float64 `json:"mvAccrued"`
	MvTotal        float64 `json:"mvTotal"`
	TargetCurrency string  `json:"targetCurrency,omitempty"`
}
