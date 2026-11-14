package persist

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

type CatalogReader interface {
	SelectInstrumentTypeID(ctx context.Context, title quik.InstrumentTypeTitle) (uint8, error)
}

type CurrencyReader interface {
	SelectCurrencies(ctx context.Context) ([]model.CurrencySource, error)
}

type QuoteReader interface {
	SelectQuotes(ctx context.Context, title quik.InstrumentTypeTitle, limit int, afterRowVersion []byte) ([]model.SourceQuote, error)
	CountBlankSecCodes(ctx context.Context, title quik.InstrumentTypeTitle) (int, error)
	LinkQuote(ctx context.Context, sourceRowVersion []byte, instrumentID int64) (bool, error)
}

type InstrumentRepository interface {
	FindInstrument(ctx context.Context, key model.InstrumentKey) (model.StoredInstrument, bool, error)
}
