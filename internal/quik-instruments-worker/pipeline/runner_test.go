package pipeline

import (
	"bytes"
	"context"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/jobs"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/persist"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	t.Run("валюта_ищется_без_создания", func(t *testing.T) {
		repo := &fakeRepo{currencies: []model.CurrencySource{{IsoCharCode: "USD"}}}
		res, err := NewRunner(nil, jobs.NewCurrencies(repo, repo, persist.NewInstruments(repo))).RunResult(context.Background(), 10)
		require.NoError(t, err)
		stats, ok := res.StatsByName("нормализация валют завершена")
		require.True(t, ok)
		require.Equal(t, 1, stats.Selected)
		require.Equal(t, 1, stats.InstrumentsUnlinked)
	})

	t.Run("котировка_связывается_с_канонической_карточкой", func(t *testing.T) {
		staleID := int64(10)
		repo := &fakeRepo{
			findInstrument: func(key model.InstrumentKey) (model.StoredInstrument, bool, error) {
				require.Equal(t, model.InstrumentKey{TypeID: 1, Ticker: "SBER", TradePointID: 1, HasTradePoint: true}, key)
				return model.StoredInstrument{InstrumentID: 42}, true, nil
			},
			selectQuotes: func(_ context.Context, title quik.InstrumentTypeTitle, _ int, after []byte) ([]model.SourceQuote, error) {
				if title != quik.InstrumentTitleSecurities || len(after) > 0 {
					return nil, nil
				}
				return []model.SourceQuote{{SourceRowVersion: []byte{1}, LinkedInstrumentID: &staleID, SecCode: "SBER", TradePointID: 1}}, nil
			},
		}
		res, err := NewRunner(nil, jobs.NewQuotes(repo, repo, persist.NewInstruments(repo), quik.InstrumentTitleSecurities, "акции")).RunResult(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, []int64{42}, repo.linkedIDs)
		stats, ok := res.StatsByName("акции")
		require.True(t, ok)
		require.Equal(t, 1, stats.QuotesLinked)
	})

	t.Run("нет_карточки_котировка_остаётся_несвязанной", func(t *testing.T) {
		repo := &fakeRepo{
			selectQuotes: func(_ context.Context, title quik.InstrumentTypeTitle, _ int, after []byte) ([]model.SourceQuote, error) {
				if title != quik.InstrumentTitleSecurities || len(after) > 0 {
					return nil, nil
				}
				return []model.SourceQuote{{SourceRowVersion: []byte{1}, SecCode: "NEW", TradePointID: 1}}, nil
			},
		}
		res, err := NewRunner(nil, jobs.NewQuotes(repo, repo, persist.NewInstruments(repo), quik.InstrumentTitleSecurities, "акции")).RunResult(context.Background(), 10)
		require.NoError(t, err)
		require.Empty(t, repo.linkedIDs)
		stats, ok := res.StatsByName("акции")
		require.True(t, ok)
		require.Equal(t, 1, stats.InstrumentsUnlinked)
	})

	t.Run("следующая_страница_обрабатывается_после_неизвестной_карточки", func(t *testing.T) {
		var lookedUp []string
		repo := &fakeRepo{
			findInstrument: func(key model.InstrumentKey) (model.StoredInstrument, bool, error) {
				lookedUp = append(lookedUp, key.Ticker)
				return model.StoredInstrument{InstrumentID: 7}, key.Ticker == "KNOWN", nil
			},
			selectQuotes: func(_ context.Context, title quik.InstrumentTypeTitle, _ int, after []byte) ([]model.SourceQuote, error) {
				if title != quik.InstrumentTitleSecurities {
					return nil, nil
				}
				if len(after) == 0 {
					return []model.SourceQuote{{SourceRowVersion: []byte{1}, SecCode: "MISSING", TradePointID: 1}}, nil
				}
				if bytes.Equal(after, []byte{1}) {
					return []model.SourceQuote{{SourceRowVersion: []byte{2}, SecCode: "KNOWN", TradePointID: 1}}, nil
				}
				return nil, nil
			},
		}
		res, err := NewRunner(nil, jobs.NewQuotes(repo, repo, persist.NewInstruments(repo), quik.InstrumentTitleSecurities, "акции")).RunResult(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, []string{"MISSING", "KNOWN"}, lookedUp)
		require.Equal(t, []int64{7}, repo.linkedIDs)
		stats, ok := res.StatsByName("акции")
		require.True(t, ok)
		require.Equal(t, 1, stats.InstrumentsUnlinked)
		require.Equal(t, 1, stats.QuotesLinked)
	})
}

type fakeRepo struct {
	currencies     []model.CurrencySource
	findInstrument func(model.InstrumentKey) (model.StoredInstrument, bool, error)
	selectQuotes   func(context.Context, quik.InstrumentTypeTitle, int, []byte) ([]model.SourceQuote, error)
	linkedIDs      []int64
}

func (r *fakeRepo) SelectInstrumentTypeID(_ context.Context, title quik.InstrumentTypeTitle) (uint8, error) {
	switch title {
	case quik.InstrumentTitleSecurities:
		return 1, nil
	case quik.InstrumentTitleBonds:
		return 2, nil
	case quik.InstrumentTitleCurrencies:
		return 3, nil
	case quik.InstrumentTitleCurrencyPairs:
		return 4, nil
	default:
		return 0, nil
	}
}

func (r *fakeRepo) SelectCurrencies(context.Context) ([]model.CurrencySource, error) {
	return r.currencies, nil
}

func (r *fakeRepo) SelectQuotes(ctx context.Context, title quik.InstrumentTypeTitle, limit int, after []byte) ([]model.SourceQuote, error) {
	if r.selectQuotes == nil {
		return nil, nil
	}
	return r.selectQuotes(ctx, title, limit, after)
}

func (*fakeRepo) CountBlankSecCodes(context.Context, quik.InstrumentTypeTitle) (int, error) {
	return 0, nil
}

func (r *fakeRepo) FindInstrument(_ context.Context, key model.InstrumentKey) (model.StoredInstrument, bool, error) {
	if r.findInstrument == nil {
		return model.StoredInstrument{}, false, nil
	}
	return r.findInstrument(key)
}

func (r *fakeRepo) LinkQuote(_ context.Context, _ []byte, instrumentID int64) (bool, error) {
	r.linkedIDs = append(r.linkedIDs, instrumentID)
	return true, nil
}
