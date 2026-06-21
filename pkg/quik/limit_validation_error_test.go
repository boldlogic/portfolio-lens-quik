package quik

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestFieldValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  FieldValidationError
		want string
	}{
		{
			name: "clientCode_длина",
			err:  FieldValidationError{Field: LimitFieldClientCode, Kind: FieldLength, Min: 1, Max: 12},
			want: "код клиента: 1-12 символов",
		},
		{
			name: "firmCode_длина",
			err:  FieldValidationError{Field: LimitFieldFirmCode, Kind: FieldLength, Min: 1, Max: 12},
			want: "код фирмы: 1-12 символов",
		},
		{
			name: "positionCode_обязателен",
			err:  FieldValidationError{Field: LimitFieldPositionCode, Kind: FieldRequired},
			want: "код позиции: обязателен",
		},
		{
			name: "positionCode_длина",
			err:  FieldValidationError{Field: LimitFieldPositionCode, Kind: FieldLength, Min: 1, Max: 4},
			want: "код позиции: 1-4 символов",
		},
		{
			name: "secCode_длина",
			err:  FieldValidationError{Field: LimitFieldSecCode, Kind: FieldLength, Min: 1, Max: 12},
			want: "код бумаги: 1-12 символов",
		},
		{
			name: "tradeAccount_обязателен",
			err:  FieldValidationError{Field: LimitFieldTradeAccount, Kind: FieldRequired},
			want: "торговый счёт: обязателен",
		},
		{
			name: "tradeAccount_длина",
			err:  FieldValidationError{Field: LimitFieldTradeAccount, Kind: FieldLength, Min: 1, Max: 12},
			want: "торговый счёт: 1-12 символов",
		},
		{
			name: "isin_длина",
			err:  FieldValidationError{Field: LimitFieldISIN, Kind: FieldLength, Min: 1, Max: 12},
			want: "isin: 1-12 символов",
		},
		{
			name: "acquisitionCurrencyCode_макс_длина",
			err:  FieldValidationError{Field: LimitFieldAcquisitionCurrencyCode, Kind: FieldMaxLength, Max: 4},
			want: "валюта приобретения: не более 4 символов",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.err.Error())
			require.ErrorIs(t, tt.err, ErrFieldValidation)
		})
	}
}

func TestValidateRuneLen(t *testing.T) {
	err := validateRuneLen(LimitFieldClientCode, "1234567890123", minClientCodeLen, maxClientCodeLen)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFieldValidation)

	var fieldErr FieldValidationError
	require.ErrorAs(t, err, &fieldErr)
	require.Equal(t, LimitFieldClientCode, fieldErr.Field)
	require.Equal(t, FieldLength, fieldErr.Kind)
}

func TestNewLimit_невалидный_clientCode(t *testing.T) {
	_, err := NewLimit(
		string(LimitTypeMoney),
		"1234567890123",
		"USD",
		new("EQTV"),
		"T0",
		nil,
		"f",
		decimal.RequireFromString("1"),
		nil,
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFieldValidation)

	var fieldErr FieldValidationError
	require.ErrorAs(t, err, &fieldErr)
	require.Equal(t, LimitFieldClientCode, fieldErr.Field)
	require.Equal(t, "код клиента: 1-12 символов", err.Error())
}

func TestNewLimit_security_без_tradeAccount(t *testing.T) {
	_, err := NewLimit(
		string(LimitTypeSecurities),
		"123",
		"LQDT",
		nil,
		"T0",
		nil,
		"f1",
		decimal.RequireFromString("1"),
		nil,
		nil,
	)
	require.Error(t, err)

	var fieldErr FieldValidationError
	require.ErrorAs(t, err, &fieldErr)
	require.Equal(t, LimitFieldTradeAccount, fieldErr.Field)
	require.Equal(t, FieldRequired, fieldErr.Kind)
	require.Equal(t, "торговый счёт: обязателен", err.Error())
}
