package jobs

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/persist"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/stats"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

const currenciesDoneMessage = "нормализация валют завершена"

type Currencies struct {
	catalog     persist.CatalogReader
	currencies  persist.CurrencyReader
	instruments *persist.Instruments
}

func NewCurrencies(catalog persist.CatalogReader, currencies persist.CurrencyReader, instruments *persist.Instruments) *Currencies {
	return &Currencies{
		catalog:     catalog,
		currencies:  currencies,
		instruments: instruments,
	}
}

func (j *Currencies) Name() string {
	return currenciesDoneMessage
}

func (j *Currencies) Run(ctx context.Context, _ int) (stats.Stats, error) {
	typeID, err := j.catalog.SelectInstrumentTypeID(ctx, quik.InstrumentTitleCurrencies)
	if err != nil {
		return stats.Stats{}, err
	}

	items, err := j.currencies.SelectCurrencies(ctx)
	if err != nil {
		return stats.Stats{}, err
	}

	drafts := make([]model.InstrumentDraft, 0, len(items))
	for _, currency := range items {
		drafts = append(drafts, currencyDraft(typeID, currency))
	}
	drafts = deduplicateInstrumentDrafts(drafts)

	result := stats.Stats{Selected: len(items)}
	for _, draft := range drafts {
		_, found, err := j.instruments.Find(ctx, draft.Key)
		if err != nil {
			return stats.Stats{}, err
		}
		if !found {
			result.InstrumentsUnlinked++
		}
	}
	return result, nil
}

func currencyDraft(typeID uint8, currency model.CurrencySource) model.InstrumentDraft {
	name := firstNotEmpty(currency.CurrencyName, currency.LatName, &currency.IsoCharCode)
	shortName := currency.IsoCharCode

	return model.InstrumentDraft{
		Key: model.InstrumentKey{
			TypeID: typeID,
			Ticker: currency.IsoCharCode,
		},
		FullName:  name,
		ShortName: &shortName,
	}
}

func deduplicateInstrumentDrafts(items []model.InstrumentDraft) []model.InstrumentDraft {
	out := make([]model.InstrumentDraft, 0, len(items))
	seen := make(map[model.InstrumentKey]int, len(items))
	for _, item := range items {
		if idx, ok := seen[item.Key]; ok {
			out[idx] = item
			continue
		}
		seen[item.Key] = len(out)
		out = append(out, item)
	}
	return out
}

func firstNotEmpty(values ...*string) *string {
	for _, value := range values {
		if value != nil && *value != "" {
			return value
		}
	}
	return nil
}
