package producer

import "github.com/boldlogic/portfolio-lens-quik/pkg/quik"

type currencyEvent struct {
	ISOCode     int16  `json:"isoCode"`
	ISOCharCode string `json:"isoCharCode"`
	Name        string `json:"name"`
	LatName     string `json:"latName"`
	MinorUnits  int32  `json:"minorUnits"`
}

func currencyToEvent(c quik.Currency) currencyEvent {
	var name string
	if c.Name() != nil {
		name = *c.Name()
	}
	return currencyEvent{
		ISOCode:     c.Numeric(),
		ISOCharCode: c.Alpha().String(),
		Name:        name,
		LatName:     c.LatName(),
		MinorUnits:  c.MinorUnits(),
	}
}
