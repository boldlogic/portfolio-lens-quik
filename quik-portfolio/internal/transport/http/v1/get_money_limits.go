package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
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

type moneyLimitsDTO struct {
	Limits     []moneyLimitDTO `json:"limits"`
	TotalCount int             `json:"totalCount"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
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

func moneyLimitsWithPaginationToResp(mls []quik.MoneyLimit, totalCount, limit, offset int) moneyLimitsDTO {
	return moneyLimitsDTO{
		Limits:     moneyLimitsToResp(mls),
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}
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

const maxClientCodeLen = 12

func normalizeClientCodes(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToUpper(strings.TrimSpace(p))
		if p != "" {
			if len(p) > maxClientCodeLen {
				return nil, fmt.Errorf("clientCode %s longer than %d chars", p, maxClientCodeLen)
			}
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

func (h *Handler) extractClientsQueryParam(r *http.Request) ([]string, error) {
	return normalizeClientCodes(r.URL.Query().Get("clientCodes"))
}

func (h *Handler) GetMoneyLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.extractDateQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	limit, offset, err := httputils.ParseListPagination(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}
	clients, err := h.extractClientsQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	mls, totalCount, err := h.service.GetMoneyLimitsWithFilters(ctx, date, limit, offset, clients)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}

		h.logger.Error("лимиты ДС: чтение HTTP", zap.Error(err), zap.Time("date", date))
		return nil, "", err
	}

	return moneyLimitsWithPaginationToResp(mls, totalCount, limit, offset), "", nil
}
