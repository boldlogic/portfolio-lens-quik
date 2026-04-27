package service

import (
	"context"
	"time"

	"github.com/boldlogic/packages/utils/dates"
)

func doRollForward(
	ctx context.Context,
	getMaxDate func(context.Context) (*time.Time, error),
	insertCopy func(context.Context, time.Time, time.Time) error,
	deleteBefore func(context.Context, time.Time) error,
) error {
	date, err := getMaxDate(ctx)
	if err != nil {
		return err
	}
	if date == nil {
		return nil
	}

	loc := time.Now().Location()
	maxDateOnly := dates.TruncateToDateIn(*date, loc)
	todayOnly := dates.Today()

	if !todayOnly.After(maxDateOnly) {
		return nil
	}

	if err := insertCopy(ctx, maxDateOnly, todayOnly); err != nil {
		return err
	}
	return deleteBefore(ctx, maxDateOnly)
}
