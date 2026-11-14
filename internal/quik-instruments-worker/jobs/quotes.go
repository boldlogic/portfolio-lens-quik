package jobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/persist"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/stats"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

type Quotes struct {
	catalog     persist.CatalogReader
	quotes      persist.QuoteReader
	instruments *persist.Instruments
	title       quik.InstrumentTypeTitle
	doneMessage string
}

func NewQuotes(
	catalog persist.CatalogReader,
	quotes persist.QuoteReader,
	instruments *persist.Instruments,
	title quik.InstrumentTypeTitle,
	doneMessage string,
) *Quotes {
	return &Quotes{
		catalog:     catalog,
		quotes:      quotes,
		instruments: instruments,
		title:       title,
		doneMessage: doneMessage,
	}
}

func (j *Quotes) Name() string {
	return j.doneMessage
}

func (j *Quotes) Run(ctx context.Context, batchSize int) (stats.Stats, error) {
	typeID, err := j.catalog.SelectInstrumentTypeID(ctx, j.title)
	if err != nil {
		return stats.Stats{}, err
	}

	blank, err := j.quotes.CountBlankSecCodes(ctx, j.title)
	if err != nil {
		return stats.Stats{}, err
	}

	result := stats.Stats{QuotesSkipped: blank}
	var after []byte
	for {
		if err := ctx.Err(); err != nil {
			return stats.Stats{}, err
		}

		quotes, err := j.quotes.SelectQuotes(ctx, j.title, batchSize, after)
		if err != nil {
			return stats.Stats{}, err
		}
		if len(quotes) == 0 {
			break
		}

		pageStats, err := j.applyPage(ctx, typeID, quotes)
		if err != nil {
			return stats.Stats{}, err
		}
		addStats(&result, pageStats)

		if len(quotes) < batchSize {
			break
		}
		after = quotes[len(quotes)-1].SourceRowVersion
	}

	return result, nil
}

func (j *Quotes) applyPage(ctx context.Context, typeID uint8, quotes []model.SourceQuote) (stats.Stats, error) {
	quoteDrafts := make([]quoteInstrumentDraft, 0, len(quotes))
	result := stats.Stats{}
	for _, quote := range quotes {
		if strings.TrimSpace(quote.SecCode) == "" {
			result.QuotesSkipped++
			continue
		}
		draft, err := quoteDraft(typeID, j.title, quote)
		if err != nil {
			return stats.Stats{}, err
		}
		quoteDrafts = append(quoteDrafts, quoteInstrumentDraft{
			Draft:  draft,
			Source: quote,
		})
	}
	result.Selected = len(quoteDrafts)

	canonicalDrafts := deduplicateQuoteDrafts(quoteDrafts)
	instruments := make(map[model.InstrumentKey]model.StoredInstrument, len(canonicalDrafts))
	for _, item := range canonicalDrafts {
		stored, found, err := j.instruments.Find(ctx, item.Draft.Key)
		if err != nil {
			return stats.Stats{}, err
		}
		if !found {
			continue
		}
		instruments[item.Draft.Key] = stored
	}

	for _, item := range quoteDrafts {
		instrument, ok := instruments[item.Draft.Key]
		if !ok {
			result.InstrumentsUnlinked++
			continue
		}

		if alreadyLinked(item.Source.LinkedInstrumentID, instrument.InstrumentID) {
			continue
		}
		linked, err := j.quotes.LinkQuote(ctx, item.Source.SourceRowVersion, instrument.InstrumentID)
		if err != nil {
			return stats.Stats{}, err
		}
		if linked {
			result.QuotesLinked++
		}
	}

	return result, nil
}

func addStats(total *stats.Stats, page stats.Stats) {
	total.Selected += page.Selected
	total.InstrumentsUnlinked += page.InstrumentsUnlinked
	total.QuotesLinked += page.QuotesLinked
	total.QuotesSkipped += page.QuotesSkipped
}

func alreadyLinked(current *int64, instrumentID int64) bool {
	return current != nil && *current == instrumentID
}

type quoteInstrumentDraft struct {
	Draft  model.InstrumentDraft
	Source model.SourceQuote
}

func quoteDraft(typeID uint8, title quik.InstrumentTypeTitle, quote model.SourceQuote) (model.InstrumentDraft, error) {
	if quote.TradePointID == 0 {
		return model.InstrumentDraft{}, fmt.Errorf("у биржевого инструмента не заполнена торговая площадка")
	}
	if quote.SecCode == "" {
		return model.InstrumentDraft{}, fmt.Errorf("у биржевого инструмента не заполнен код")
	}

	draft := model.InstrumentDraft{
		Key: model.InstrumentKey{
			TypeID:        typeID,
			Ticker:        quote.SecCode,
			TradePointID:  quote.TradePointID,
			HasTradePoint: true,
		},
		RegistrationNumber: quote.RegistrationNumber,
		FullName:           quote.FullName,
		ShortName:          quote.ShortName,
		ISIN:               quote.ISIN,
	}

	if title == quik.InstrumentTitleBonds {
		draft.FaceValue = quote.FaceValue
		draft.MaturityDate = quote.MaturityDate
		draft.CouponDuration = quote.CouponDuration
	}

	return draft, nil
}

func deduplicateQuoteDrafts(items []quoteInstrumentDraft) []quoteInstrumentDraft {
	out := make([]quoteInstrumentDraft, 0, len(items))
	seen := make(map[model.InstrumentKey]int, len(items))
	for _, item := range items {
		if idx, ok := seen[item.Draft.Key]; ok {
			if model.RowVersionAfter(item.Source.SourceRowVersion, out[idx].Source.SourceRowVersion) {
				out[idx] = item
			}
			continue
		}
		seen[item.Draft.Key] = len(out)
		out = append(out, item)
	}
	return out
}
