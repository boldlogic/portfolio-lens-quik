package models

import "errors"

type SettleCode string

const (
	SettleCodeT0 = "T0"
	SettleCodeT1 = "T1"
	SettleCodeT2 = "T2"
	SettleCodeTx = "Tx"
)

var ErrWrongSettleCode = errors.New("settleCode должен быть T0, T1, T2 или Tx")

func (s SettleCode) Validate() error {
	switch s {
	case SettleCodeT0, SettleCodeT1, SettleCodeT2, SettleCodeTx:
		return nil
	default:
		return ErrWrongSettleCode
	}
}

func (s SettleCode) String() string {
	switch s {
	case SettleCodeT0:
		return "T0"
	case SettleCodeT1:
		return "T1"
	case SettleCodeT2:
		return "T2"
	case SettleCodeTx:
		return "Tx"
	default:
		return "unknown"
	}
}
