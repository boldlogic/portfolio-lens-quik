package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/shopspring/decimal"
)

func (h *Handler) getMoneyLimits(r *http.Request) (any, string, error) {
	q, err := parseLimitsQueryParams(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	mls, totalCount, err := h.service.GetMoneyLimitsWithFilters(
		r.Context(), q.LoadDate, q.Limit, q.Offset, q.ClientCodes, q.IncludeTotalCount,
	)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return moneyLimitsToResponseDTO(mls, q.Limit, q.Offset, totalCount, q.IncludeTotalCount), "", nil
}

type moneyLimitsResponseDTO struct {
	Limits     []moneyLimitDTO `json:"limits"`
	TotalCount *uint64         `json:"totalCount,omitempty"`
	Limit      uint32          `json:"limit"`
	Offset     uint64          `json:"offset"`
}

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

func moneyLimitsToResponseDTO(mls []quik.MoneyLimit, limit uint32, offset uint64, totalCount *uint64, includeTotalCount bool) moneyLimitsResponseDTO {

	if includeTotalCount && totalCount == nil {
		var z uint64 = 0
		totalCount = new(z)
	}

	out := moneyLimitsResponseDTO{
		Limits:     moneyLimitsToDTO(mls),
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}
	return out
}

func moneyLimitsToDTO(mls []quik.MoneyLimit) []moneyLimitDTO {
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
		Currency:     currencyCodeForDTO(ml.CurrencyCode),
		PositionCode: ml.PositionCode,
		SettleCode:   string(ml.SettleCode),
		FirmCode:     ml.FirmCode,
		FirmName:     ml.FirmName,
		Balance:      ml.Balance,
	}
}
