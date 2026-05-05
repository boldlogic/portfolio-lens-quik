package v1

import (
	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type moneyLimitDTO struct {
	LoadDate     string          `json:"loadDate"`
	SourceDate   string          `json:"sourceDate"`
	ClientCode   string          `json:"clientCode"`
	Currency     string          `json:"currency"`
	PositionCode string          `json:"positionCode"`
	SettleCode   string          `json:"settleCode"`
	FirmCode     string          `json:"firmCode"`
	FirmName     string          `json:"firmName"`
	Balance      decimal.Decimal `json:"balance"`
}

func moneyLimitsToResp(mls []quik.MoneyLimit) []moneyLimitDTO {
	if len(mls) == 0 {
		return []moneyLimitDTO{}
	}

	resp := make([]moneyLimitDTO, 0, len(mls))
	for _, ml := range mls {
		resp = append(resp, moneyLimitToDTO(ml))
	}
	return resp
}

func moneyLimitToDTO(ml quik.MoneyLimit) moneyLimitDTO {
	return moneyLimitDTO{
		LoadDate:     ml.LoadDate.Format(dates.ISODateFormat),
		SourceDate:   ml.SourceDate.Format(dates.ISODateFormat),
		ClientCode:   ml.ClientCode,
		Currency:     ml.Currency,
		PositionCode: ml.PositionCode,
		SettleCode:   string(ml.SettleCode),
		FirmCode:     ml.FirmCode,
		FirmName:     ml.FirmName,
		Balance:      ml.Balance,
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
	return out
}
