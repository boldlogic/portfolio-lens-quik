package v1

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/boldlogic/packages/transport/httputils"
	intmodels "github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/shopspring/decimal"
)

func processRowCSV(line int, idx map[int]string, row []string, res *[]intmodels.LimitLine) error {
	var l limitCSV

	for i, v := range row {
		if line == 1 {
			idx[i] = v
		}
		switch idx[i] {
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
		case "acquisition_currency":
			l.AcquisitionCurrencyCode = v
		}

	}

	if line > 1 {
		balance, err := decimal.NewFromString(l.Balance)
		if err != nil {
			return fmt.Errorf("некорректный баланс в строке %d", line)
		}
		line := intmodels.LimitLine{
			LimitInput: intmodels.LimitInput{
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
			Line: uint(line),
		}
		*res = append(*res, line)
	}
	return nil

}

func (h *Handler) upload(r *http.Request) (any, string, error) {
	ctx := r.Context()
	mr, err := r.MultipartReader()
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err.Error(), httputils.ErrUnsupportedMediaType
	}

	res := make([]intmodels.LimitLine, 0, 100)

outer:
	for {
		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()

		default:
			p, err := mr.NextPart()
			if errors.Is(err, io.EOF) {
				break outer
			}

			if err != nil {
				return nil, err.Error(), models.ErrValidation
			}

			ct := p.Header.Get("Content-Type")

			if !strings.Contains(ct, "csv") {
				return nil, fmt.Sprintf("Поддерживается только csv, получен %s", ct), httputils.ErrUnsupportedMediaType
			}

			cr := csv.NewReader(p)
			cr.Comma = ';'

			idx := make(map[int]string, 9)
			line := 0

		inner:
			for {
				select {
				case <-ctx.Done():
					return nil, "", ctx.Err()

				default:
					line++
					row, err := cr.Read()
					if errors.Is(err, io.EOF) {
						break inner
					}
					if err != nil {
						return nil, err.Error(), models.ErrValidation
					}
					err = processRowCSV(line, idx, row, &res)
					if err != nil {
						return nil, err.Error(), models.ErrValidation
					}
				}
			}
		}

	}

	err = h.service.UpsertLimits(r.Context(), res)
	if err != nil {
		return nil, err.Error(), err
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
