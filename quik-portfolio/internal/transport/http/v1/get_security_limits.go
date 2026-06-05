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

func (h *Handler) getSecurityLimits(r *http.Request) (any, string, error) {
	q, err := parseLimitsListQuery(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, totalCount, err := h.service.GetSecurityLimitsWithFilters(
		r.Context(), q.Date, q.Limit, q.Offset, q.ClientCodes, q.IncludeTotalCount,
	)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return securityLimitsWithPaginationToResp(sls, q.Limit, q.Offset, totalCount, q.IncludeTotalCount), "", nil
}

type securityLimitsDTO struct {
	Limits     []securityLimitDTO `json:"limits"`
	TotalCount *uint64            `json:"totalCount,omitempty"`
	Limit      uint32             `json:"limit"`
	Offset     uint64             `json:"offset"`
}

type securityLimitDTO struct {
	LoadDate       string          `json:"loadDate"`
	SourceDate     string          `json:"sourceDate"`
	ClientCode     string          `json:"clientCode"`
	Ticker         string          `json:"ticker"`
	TradeAccount   string          `json:"tradeAccount"`
	SettleCode     string          `json:"settleCode"`
	FirmCode       string          `json:"firmCode"`
	FirmName       string          `json:"firmName"`
	Balance        decimal.Decimal `json:"balance"`
	AcquisitionCcy string          `json:"acquisitionCcy"`
	ISIN           string          `json:"isin,omitempty"`
	ShortName      string          `json:"shortName,omitempty"`
}

func securityLimitsWithPaginationToResp(sls []quik.SecurityLimit, limit uint32, offset uint64, totalCount *uint64, includeTotalCount bool) securityLimitsDTO {
	if includeTotalCount && totalCount == nil {
		var z uint64 = 0
		totalCount = new(z)
	}
	return securityLimitsDTO{
		Limits:     securityLimitsToResp(sls),
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}
}

func securityLimitsToResp(sls []quik.SecurityLimit) []securityLimitDTO {
	res := make([]securityLimitDTO, 0, len(sls))
	for _, sl := range sls {
		res = append(res, securityLimitToDTO(sl))
	}
	return res
}

func securityLimitToDTO(sl quik.SecurityLimit) securityLimitDTO {
	var out securityLimitDTO
	out.LoadDate = sl.LoadDate.Format(dates.ISODateFormat)
	out.SourceDate = sl.SourceDate.Format(dates.ISODateFormat)
	out.ClientCode = sl.ClientCode
	out.Ticker = sl.Ticker
	out.TradeAccount = sl.TradeAccount
	out.SettleCode = string(sl.SettleCode)
	out.FirmCode = sl.FirmCode
	out.FirmName = sl.FirmName
	out.Balance = sl.Balance
	out.AcquisitionCcy = sl.AcquisitionCcy

	out.ISIN = sl.ISIN

	if sn := strings.TrimSpace(sl.ShortName); sn != "" {
		out.ShortName = sn
	}
	return out
}
