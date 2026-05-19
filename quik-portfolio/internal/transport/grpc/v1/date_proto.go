package v1

import (
	"fmt"
	"time"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	datepb "google.golang.org/genproto/googleapis/type/date"
)

func timeToProtoDate(t time.Time) *datepb.Date {
	y, m, d := t.Date()
	return &datepb.Date{
		Year:  int32(y),
		Month: int32(m),
		Day:   int32(d),
	}
}

func protoDateToTime(d *datepb.Date) (time.Time, error) {
	if d == nil {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}

	y := int(d.Year)
	m := int(d.Month)
	day := int(d.Day)

	if d.Year < 2000 || d.Year > 2999 || d.Month < 1 || d.Month > 12 || d.Day < 1 || d.Day > 31 {
		return time.Time{}, fmt.Errorf("%w: некорректная дата", md.ErrValidation)
	}

	t := time.Date(y, time.Month(m), day, 0, 0, 0, 0, time.Local)

	return t, nil
}
