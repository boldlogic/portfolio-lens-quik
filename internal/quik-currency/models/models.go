package models

type CurrencyFromCrossrates struct {
	IsoCharCode string
	Name        string
}

type Currency struct {
	CharCode   string
	Name       *string
	Numeric    int16
	LatName    string
	MinorUnits int32
}
