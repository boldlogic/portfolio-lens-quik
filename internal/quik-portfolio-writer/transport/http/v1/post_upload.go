package v1

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strings"

	"github.com/boldlogic/packages/transport/httputils"
	intmodels "github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) upload(r *http.Request) (any, string, error) {

	err := r.ParseMultipartForm(5 << 20)
	defer r.MultipartForm.RemoveAll()

	if err != nil {
		h.logger.Error("", zap.Error(err))

		return nil, "", err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("", zap.Error(err))

		return nil, "", err
	}
	if !strings.Contains(header.Header.Get("Content-Type"), "csv") {
		return nil, "", httputils.ErrUnsupportedMediaType
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		h.logger.Error(err.Error())

		return nil, "", err
	}

	idx := make(map[int]string, 9)

	var res []intmodels.LimitLine
	for i, str := range records {
		var l limitCSV
		for j, v := range str {
			if i == 0 {
				idx[j] = v
			}

			switch idx[j] {
			case "limit_type":
				l.Type = v
			case "client_code":
				l.ClientCode = v
			case "ticker":
				l.Ticker = v
			case "settle_code":
				l.SettleCode = v
			case "firm_code":
				l.FirmCode = v
			case "balance":
				l.Balance = v
			case "isin":
				l.ISIN = v
			case "position_code":
				l.PositionCode = v
			case "trade_account":
				l.TradeAccount = v
			}

		}
		if i > 1 {
			balance, err := decimal.NewFromString(l.Balance)
			if err != nil {
				h.logger.Error(err.Error())
				return nil, "", err
			}
			line := intmodels.LimitLine{
				Limit: intmodels.Limit{
					Type:                    l.Type,
					ClientCode:              l.ClientCode,
					Ticker:                  l.Ticker,
					PositionCode:            &l.PositionCode,
					SettleCode:              l.SettleCode,
					TradeAccount:            &l.TradeAccount,
					FirmCode:                l.FirmCode,
					Balance:                 balance,
					AcquisitionCurrencyCode: &l.AcquisitionCurrencyCode,
					ISIN:                    &l.ISIN,
				},
				Line: uint(i),
			}

			res = append(res, line)
		}

	}
	err = h.service.UpsertLimits(r.Context(), res)
	if err != nil {
		if errors.Is(err, models.ErrBusinessValidation) || errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return nil, "", nil
}

type limitCSV struct {
	Type                    string
	ClientCode              string
	Ticker                  string
	PositionCode            string
	SettleCode              string
	TradeAccount            string
	FirmCode                string
	Balance                 string
	AcquisitionCurrencyCode string
	ISIN                    string
}
