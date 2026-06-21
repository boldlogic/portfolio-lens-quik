package service

import (
	"fmt"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func validateLimitsContract(date time.Time, clientCodes []string) ([]string, error) {
	if err := checkLimitDate(date); err != nil {
		return nil, err
	}
	dedublicated, err := deduplicateClientCodes(clientCodes)
	if err != nil {
		return nil, err
	}
	return dedublicated, nil
}

func checkLimitDate(date time.Time) error {
	min := dates.Today().AddDate(0, 0, -1)
	max := dates.Today().AddDate(0, 0, 1)
	if date.Before(min) || date.After(max) {
		return fmt.Errorf("%w: дата должна быть в диапазоне от %s до %s",
			models.ErrBusinessValidation,
			min.Format(dates.ISODateFormat),
			max.Format(dates.ISODateFormat))
	}

	return nil
}

const maxClientCodesCount = 10

func deduplicateClientCodes(clientCodes []string) ([]string, error) {
	if len(clientCodes) == 0 {
		return nil, nil
	}

	dedup := make(map[string]struct{}, maxClientCodesCount)
	out := make([]string, 0, maxClientCodesCount)

	unique := 0
	for _, code := range clientCodes {
		trimmed, err := quik.ParseClientCode(code)
		if err != nil {
			return nil, fmt.Errorf("%w: код клиента %s", models.ErrBusinessValidation, code)
		}
		if _, ok := dedup[trimmed]; ok {
			continue
		}

		if unique++; unique > maxClientCodesCount {
			return nil, fmt.Errorf("%w: слишком много клиентских кодов, ограничение %d уникальных", models.ErrBusinessValidation, maxClientCodesCount)

		}
		out = append(out, trimmed)
		dedup[trimmed] = struct{}{}
	}

	return out, nil

}
