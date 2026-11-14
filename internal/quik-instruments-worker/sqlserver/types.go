package sqlserver

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

func stringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func int16Ptr(v sql.NullInt16) *int16 {
	if !v.Valid {
		return nil
	}
	return &v.Int16
}

func int32Ptr(v sql.NullInt32) *int32 {
	if !v.Valid {
		return nil
	}
	return &v.Int32
}

func timePtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Time
}

func decimalPtr(v *decimal.Decimal) *decimal.Decimal {
	if v == nil {
		return nil
	}
	out := *v
	return &out
}
