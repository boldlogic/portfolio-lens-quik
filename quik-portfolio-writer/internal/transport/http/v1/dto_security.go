package v1

import (
	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

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
	return out
}
