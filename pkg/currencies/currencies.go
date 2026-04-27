package currencies

type Currency struct {
	ISOCode     int16
	ISOCharCode string
	Name        *string
	LatName     string
	MinorUnits  int32
}
