package quik

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type LimitField string

const (
	LimitFieldClientCode              LimitField = "clientCode"
	LimitFieldFirmCode                LimitField = "firmCode"
	LimitFieldPositionCode            LimitField = "positionCode"
	LimitFieldSecCode                 LimitField = "secCode"
	LimitFieldTradeAccount            LimitField = "tradeAccount"
	LimitFieldISIN                    LimitField = "isin"
	LimitFieldAcquisitionCurrencyCode LimitField = "acquisitionCurrencyCode"
)

type FieldValidationKind int

const (
	FieldRequired FieldValidationKind = iota
	FieldLength
	FieldMaxLength
)

type FieldValidationError struct {
	Field LimitField
	Kind  FieldValidationKind
	Min   int
	Max   int
}

var ErrFieldValidation = errors.New("ошибка поля лимита")

func (e FieldValidationError) Error() string {
	label := limitFieldLabel(e.Field)
	switch e.Kind {
	case FieldRequired:
		return fmt.Sprintf("%s: обязателен", label)
	case FieldLength:
		return fmt.Sprintf("%s: %d-%d символов", label, e.Min, e.Max)
	case FieldMaxLength:
		return fmt.Sprintf("%s: не более %d символов", label, e.Max)
	default:
		return fmt.Sprintf("%s: некорректное значение", label)
	}
}

func (e FieldValidationError) Is(target error) bool {
	return target == ErrFieldValidation
}

func limitFieldLabel(field LimitField) string {
	switch field {
	case LimitFieldClientCode:
		return "код клиента"
	case LimitFieldFirmCode:
		return "код фирмы"
	case LimitFieldPositionCode:
		return "код позиции"
	case LimitFieldSecCode:
		return "код бумаги"
	case LimitFieldTradeAccount:
		return "торговый счёт"
	case LimitFieldISIN:
		return "isin"
	case LimitFieldAcquisitionCurrencyCode:
		return "валюта приобретения"
	default:
		return string(field)
	}
}

func fieldValidationError(field LimitField, kind FieldValidationKind, min, max int) error {
	return FieldValidationError{
		Field: field,
		Kind:  kind,
		Min:   min,
		Max:   max,
	}
}

func validateRuneLen(field LimitField, s string, min, max int) error {
	n := utf8.RuneCountInString(s)
	if n < min || n > max {
		return fieldValidationError(field, FieldLength, min, max)
	}
	return nil
}

func validateRequiredPtr(field LimitField, p *string) error {
	if p == nil {
		return fieldValidationError(field, FieldRequired, 0, 0)
	}
	return nil
}

func validateMaxRuneLen(field LimitField, s string, max int) error {
	if utf8.RuneCountInString(s) > max {
		return fieldValidationError(field, FieldMaxLength, 0, max)
	}
	return nil
}

func validateOptionalMaxRuneLen(field LimitField, p *string, max int) error {
	if p == nil || strings.TrimSpace(*p) == "" {
		return nil
	}
	return validateMaxRuneLen(field, *p, max)
}
