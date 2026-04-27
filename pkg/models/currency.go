package models

import "time"

type Currency struct {
	ISOCode     int16
	ISOCharCode string
	Name        *string
	LatName     string
	MinorUnits  int32
	ExtSystemId *ExternalSystemID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
