package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	qmodels "github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
)

func (h *Handler) GetSecurityLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, err := h.service.GetSecurityLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}

		return nil, "", err
	}
	return securityLimitsToResp(sls), "", nil
}

func securityLimitsToResp(sls []qmodels.SecurityLimit) []securityLimitDTO {

	if len(sls) == 0 {
		return []securityLimitDTO{}
	}

	res := make([]securityLimitDTO, 0, len(sls))
	for _, sl := range sls {
		res = append(res, securityLimitToDTO(sl))
	}
	return res
}

func securityLimitToDTO(sl qmodels.SecurityLimit) securityLimitDTO {
	var out securityLimitDTO
	out.LoadDate = sl.LoadDate.Format(dates.ISODateFormat)
	out.SourceDate = sl.SourceDate.Format(dates.ISODateFormat)
	out.ClientCode = sl.ClientCode
	out.Ticker = sl.Ticker
	out.TradeAccount = sl.TradeAccount
	out.SettleCode = string(sl.SettleCode)
	out.FirmCode = sl.FirmCode
	out.FirmName = sl.FirmName
	out.Balance = sl.Balance.InexactFloat64()
	out.AcquisitionCcy = sl.AcquisitionCcy

	if sl.ISIN != nil {
		out.ISIN = *sl.ISIN
	}
	return out
}

type securityLimitDTO struct {
	LoadDate       string  `json:"loadDate"`
	SourceDate     string  `json:"sourceDate"`
	ClientCode     string  `json:"clientCode"`
	Ticker         string  `json:"ticker"`
	TradeAccount   string  `json:"tradeAccount"`
	SettleCode     string  `json:"settleCode"`
	FirmCode       string  `json:"firmCode"`
	FirmName       string  `json:"firmName"`
	Balance        float64 `json:"balance"`
	AcquisitionCcy string  `json:"acquisitionCcy"`
	ISIN           string  `json:"isin,omitempty"`
}
