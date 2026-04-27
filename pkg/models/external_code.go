package models

type ExternalCode struct {
	ExtId            int32
	ExternalSystemId ExternalSystemID
	Code             string
	Type             ExternalCodeType
	IntId            int64
}

type ExternalCodeType uint8

const (
	ExCodeTypeCurrency   ExternalCodeType = 1
	ExCodeTypeInstrument ExternalCodeType = 2
)
