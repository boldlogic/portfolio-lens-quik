package v1

import (
	"strings"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type securityLimitsDTO struct {
	Limits     []securityLimitDTO `json:"limits"`
	TotalCount int                `json:"totalCount"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
}

func securityLimitsWithPaginationToResp(sls []quik.SecurityLimit, totalCount, limit, offset int) securityLimitsDTO {
	return securityLimitsDTO{
		Limits:     securityLimitsToResp(sls),
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}
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

func securityLimitsToResp(sls []quik.SecurityLimit) []securityLimitDTO {

	if len(sls) == 0 {
		return []securityLimitDTO{}
	}

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

	if sl.ISIN != nil {
		out.ISIN = *sl.ISIN
	}
	if sn := strings.TrimSpace(sl.ShortName); sn != "" {
		out.ShortName = sn
	}
	return out
}
