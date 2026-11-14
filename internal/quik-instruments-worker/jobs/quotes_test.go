package jobs

import (
	"testing"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCurrencyDraft(t *testing.T) {
	name := "Доллар США"
	latName := "US Dollar"

	draft := currencyDraft(3, model.CurrencySource{
		IsoCharCode:  "USD",
		CurrencyName: &name,
		LatName:      &latName,
	})

	require.Equal(t, uint8(3), draft.Key.TypeID)
	require.Equal(t, "USD", draft.Key.Ticker)
	require.False(t, draft.Key.HasTradePoint)
	require.Equal(t, &name, draft.FullName)
	require.Nil(t, draft.ISIN)
	require.Nil(t, draft.RegistrationNumber)
	require.Nil(t, draft.FaceValue)
	require.Nil(t, draft.MaturityDate)
	require.Nil(t, draft.CouponDuration)
}

func TestQuoteDraft(t *testing.T) {
	t.Run("акция_без_bond_полей", func(t *testing.T) {
		faceValue := decimal.NewFromInt(1000)
		maturityDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		couponDuration := int32(182)

		draft, err := quoteDraft(1, quik.InstrumentTitleSecurities, model.SourceQuote{
			SecCode:        "SBER",
			TradePointID:   1,
			FaceValue:      &faceValue,
			MaturityDate:   &maturityDate,
			CouponDuration: &couponDuration,
		})

		require.NoError(t, err)
		require.Equal(t, "SBER", draft.Key.Ticker)
		require.Nil(t, draft.FaceValue)
		require.Nil(t, draft.MaturityDate)
		require.Nil(t, draft.CouponDuration)
	})

	t.Run("облигация_с_bond_полями", func(t *testing.T) {
		faceValue := decimal.NewFromInt(1000)
		maturityDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		couponDuration := int32(182)

		draft, err := quoteDraft(2, quik.InstrumentTitleBonds, model.SourceQuote{
			SecCode:        "RU000A",
			TradePointID:   1,
			FaceValue:      &faceValue,
			MaturityDate:   &maturityDate,
			CouponDuration: &couponDuration,
		})

		require.NoError(t, err)
		require.Equal(t, &faceValue, draft.FaceValue)
		require.Equal(t, &maturityDate, draft.MaturityDate)
		require.Equal(t, &couponDuration, draft.CouponDuration)
	})

	t.Run("валютная_пара_без_bond_полей", func(t *testing.T) {
		faceValue := decimal.NewFromInt(1000)
		maturityDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		couponDuration := int32(182)

		draft, err := quoteDraft(4, quik.InstrumentTitleCurrencyPairs, model.SourceQuote{
			SecCode:        "USD000UTSTOM",
			TradePointID:   1,
			FaceValue:      &faceValue,
			MaturityDate:   &maturityDate,
			CouponDuration: &couponDuration,
		})

		require.NoError(t, err)
		require.Equal(t, "USD000UTSTOM", draft.Key.Ticker)
		require.Nil(t, draft.FaceValue)
		require.Nil(t, draft.MaturityDate)
		require.Nil(t, draft.CouponDuration)
	})
}

func TestDeduplicateQuoteDrafts(t *testing.T) {
	oldQuote := model.SourceQuote{
		SourceRowVersion: []byte{0, 1},
		SecCode:          "SBER",
		TradePointID:     1,
	}
	newQuote := model.SourceQuote{
		SourceRowVersion: []byte{0, 2},
		SecCode:          "SBER",
		TradePointID:     1,
	}
	oldDraft, err := quoteDraft(1, quik.InstrumentTitleSecurities, oldQuote)
	require.NoError(t, err)
	newDraft, err := quoteDraft(1, quik.InstrumentTitleSecurities, newQuote)
	require.NoError(t, err)

	out := deduplicateQuoteDrafts([]quoteInstrumentDraft{
		{Draft: oldDraft, Source: oldQuote},
		{Draft: newDraft, Source: newQuote},
	})

	require.Len(t, out, 1)
	require.Equal(t, []byte{0, 2}, out[0].Source.SourceRowVersion)
}

func TestDeduplicateQuoteDrafts_разные_sec_code_не_склеиваются(t *testing.T) {
	tomQuote := model.SourceQuote{
		SourceRowVersion: []byte{0, 1},
		SecCode:          "USD000UTSTOM",
		TradePointID:     1,
	}
	todQuote := model.SourceQuote{
		SourceRowVersion: []byte{0, 2},
		SecCode:          "USD000UTSTOD",
		TradePointID:     1,
	}
	tomDraft, err := quoteDraft(4, quik.InstrumentTitleCurrencyPairs, tomQuote)
	require.NoError(t, err)
	todDraft, err := quoteDraft(4, quik.InstrumentTitleCurrencyPairs, todQuote)
	require.NoError(t, err)

	out := deduplicateQuoteDrafts([]quoteInstrumentDraft{
		{Draft: tomDraft, Source: tomQuote},
		{Draft: todDraft, Source: todQuote},
	})

	require.Len(t, out, 2)
}
