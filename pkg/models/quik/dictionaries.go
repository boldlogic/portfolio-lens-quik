package quik

import "github.com/boldlogic/portfolio-lens-quik/pkg/models"

type InstrumentType struct {
	Id    uint8
	Title string
}

type InstrumentSubType struct {
	SubTypeId uint8
	Title     string
	TypeId    uint8
}

type Firm struct {
	Id   uint8
	Code string
	Name string
}

type Board struct {
	Id           uint8
	Code         string
	Name         string
	IsTraded     bool
	TradePointId *uint8
	TradePoint   *models.TradePoint
}
