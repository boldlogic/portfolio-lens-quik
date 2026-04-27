package models

type ExternalSystem struct {
	Id     ExternalSystemID
	System string
}
type ExternalSystemID uint8

const (
	CBRSystem  ExternalSystemID = 1
	QuikSystem ExternalSystemID = 2
)
