package quik

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestNewLimit(t *testing.T) {
	tests := []struct {
		name                    string
		limitType               string
		clientCode              string
		ticker                  string
		positionCode            *string
		settleCode              string
		tradeAccount            *string
		firmCode                string
		balance                 decimal.Decimal
		acquisitionCurrencyCode *string
		ISIN                    *string
		want                    Limit
		wantErr                 error
	}{
		{
			name:         "money_limit_корректный_usd",
			limitType:    string(LimitTypeMoney),
			clientCode:   "AABBCC12",
			ticker:       "USD",
			positionCode: new("EQTV"),
			settleCode:   "T0",
			firmCode:     "NC0058900000",
			balance:      decimal.RequireFromString("10.00"),
			want: Limit{
				limitType:               LimitTypeMoney,
				clientCode:              "AABBCC12",
				ticker:                  "USD",
				positionCode:            new("EQTV"),
				settleCode:              "T0",
				firmCode:                "NC0058900000",
				balance:                 decimal.RequireFromString("10.00"),
				tradeAccount:            nil,
				acquisitionCurrencyCode: nil,
				isin:                    nil,
			},
		},
		{
			name:         "money_limit_корректный_SUR",
			limitType:    string(LimitTypeMoney),
			clientCode:   "123",
			ticker:       "SUR",
			positionCode: new("EQTV"),
			settleCode:   "T1",
			firmCode:     "f",
			balance:      decimal.RequireFromString("0.00"),
			want: Limit{
				limitType:    LimitTypeMoney,
				clientCode:   "123",
				ticker:       "SUR",
				positionCode: new("EQTV"),
				settleCode:   "T1",
				firmCode:     "f",
				balance:      decimal.RequireFromString("0.00"),
			},
		},
		{
			name:         "money_limit_GLD",
			limitType:    string(LimitTypeMoney),
			clientCode:   "123",
			ticker:       "GLD",
			positionCode: new("EQTV"),
			settleCode:   "T1",
			firmCode:     "f",
			balance:      decimal.RequireFromString("-10.00"),
			want: Limit{
				limitType:    LimitTypeMoney,
				clientCode:   "123",
				ticker:       "GLD",
				positionCode: new("EQTV"),
				settleCode:   "T1",
				firmCode:     "f",
				balance:      decimal.RequireFromString("-10.00"),
			},
		},
		{
			name:         "sec_limit_неполный",
			limitType:    string(LimitTypeSecurities),
			clientCode:   "AABBCC12",
			ticker:       "LQDT",
			tradeAccount: new("L01-00000F01"),
			settleCode:   "T0",
			firmCode:     "f1",
			balance:      decimal.RequireFromString("10.00"),
			want: Limit{
				limitType:               LimitTypeSecurities,
				clientCode:              "AABBCC12",
				ticker:                  "LQDT",
				tradeAccount:            new("L01-00000F01"),
				settleCode:              "T0",
				firmCode:                "f1",
				balance:                 decimal.RequireFromString("10.00"),
				positionCode:            nil,
				acquisitionCurrencyCode: nil,
				isin:                    nil,
			},
		},
		{
			name:                    "sec_limit_полный",
			limitType:               string(LimitTypeSecurities),
			clientCode:              "AABBCC12",
			ticker:                  "BCSR",
			tradeAccount:            new("L01-00000F01"),
			settleCode:              "T2",
			firmCode:                "f1",
			balance:                 decimal.RequireFromString("10.00"),
			acquisitionCurrencyCode: new("SUR"),
			ISIN:                    new("RU000A10A0N6"),
			want: Limit{
				limitType:               LimitTypeSecurities,
				clientCode:              "AABBCC12",
				ticker:                  "BCSR",
				tradeAccount:            new("L01-00000F01"),
				settleCode:              "T2",
				firmCode:                "f1",
				balance:                 decimal.RequireFromString("10.00"),
				acquisitionCurrencyCode: new("SUR"),
				isin:                    new("RU000A10A0N6"),
			},
		},
		{
			name:                    "sec_limit_облигация",
			limitType:               string(LimitTypeSecurities),
			clientCode:              "AABBCC12",
			ticker:                  "SU26248RMFS3",
			tradeAccount:            new("L01-00000F01"),
			settleCode:              "T1",
			firmCode:                "f1",
			balance:                 decimal.RequireFromString("10.00"),
			acquisitionCurrencyCode: new("%"),
			ISIN:                    new("SU26248RMFS3"),
			want: Limit{
				limitType:               LimitTypeSecurities,
				clientCode:              "AABBCC12",
				ticker:                  "SU26248RMFS3",
				tradeAccount:            new("L01-00000F01"),
				settleCode:              "T1",
				firmCode:                "f1",
				balance:                 decimal.RequireFromString("10.00"),
				acquisitionCurrencyCode: new("%"),
				isin:                    new("SU26248RMFS3"),
			},
		},
		{
			name:         "otc",
			limitType:    string(LimitTypeSecuritiesOtc),
			clientCode:   "AABBCC12",
			ticker:       "LQDT",
			tradeAccount: new("L01-00000F01"),
			settleCode:   "T0",
			firmCode:     "f1",
			balance:      decimal.RequireFromString("10.00"),
			want: Limit{
				limitType:               LimitTypeSecuritiesOtc,
				clientCode:              "AABBCC12",
				ticker:                  "LQDT",
				tradeAccount:            new("L01-00000F01"),
				settleCode:              "T0",
				firmCode:                "f1",
				balance:                 decimal.RequireFromString("10.00"),
				positionCode:            nil,
				acquisitionCurrencyCode: nil,
				isin:                    nil,
			},
		},
		{
			name:       "otc_неполный",
			limitType:  string(LimitTypeSecuritiesOtc),
			clientCode: "AABBCC12",
			ticker:     "LQDT",
			firmCode:   "f1",
			balance:    decimal.RequireFromString("10.00"),
			want: Limit{
				limitType:               LimitTypeSecuritiesOtc,
				clientCode:              "AABBCC12",
				ticker:                  "LQDT",
				tradeAccount:            new("OTC"),
				settleCode:              SettleCodeTx,
				firmCode:                "f1",
				balance:                 decimal.RequireFromString("10.00"),
				positionCode:            nil,
				acquisitionCurrencyCode: nil,
				isin:                    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewLimit(tt.limitType, tt.clientCode, tt.ticker, tt.positionCode, tt.settleCode, tt.tradeAccount, tt.firmCode, tt.balance, tt.acquisitionCurrencyCode, tt.ISIN)
			if tt.wantErr != nil {
				require.ErrorAs(t, gotErr, tt.wantErr)
				require.Empty(t, got)
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)

		})
	}
}
