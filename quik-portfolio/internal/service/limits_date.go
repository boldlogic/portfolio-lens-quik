package service

import (
	"fmt"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func checkLimitDate(date time.Time, minAvailable time.Time) error {
	loc := time.Now().Location()
	today := dates.Today()
	dateTrunc := dates.TruncateToDateIn(date, loc)
	minTrunc := dates.TruncateToDateIn(minAvailable, loc)

	if dateTrunc.Before(minTrunc) || dateTrunc.After(today) {
		return fmt.Errorf("%w: дата должна быть в диапазоне от %s до %s",
			models.ErrBusinessValidation,
			minTrunc.Format(dates.ISODateFormat),
			today.Format(dates.ISODateFormat),
		)
	}
	return nil
}

func minRollForwardDate(maxDate *time.Time) time.Time {
	if maxDate == nil {
		return dates.Today()
	}
	return *maxDate
}
