package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
)

const selectInstrumentSQL = `
SELECT instrument_id, type_id, trade_point_id, ticker
FROM ref.instruments
WHERE type_id = @p1
	AND ticker = @p2
	AND ((@p3 IS NULL AND trade_point_id IS NULL) OR trade_point_id = @p3);`

func (r *Repository) FindInstrument(ctx context.Context, key model.InstrumentKey) (model.StoredInstrument, bool, error) {
	var out model.StoredInstrument
	var tradePointID sql.NullByte
	err := r.Db.QueryRowContext(ctx, selectInstrumentSQL, key.TypeID, key.Ticker, tradePointIDParam(key)).Scan(
		&out.InstrumentID,
		&out.Draft.Key.TypeID,
		&tradePointID,
		&out.Draft.Key.Ticker,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.StoredInstrument{}, false, nil
		}
		return model.StoredInstrument{}, false, fmt.Errorf("ошибка чтения инструмента: %w", err)
	}
	if tradePointID.Valid {
		out.Draft.Key.TradePointID = tradePointID.Byte
		out.Draft.Key.HasTradePoint = true
	}
	return out, true, nil
}

func tradePointIDParam(key model.InstrumentKey) any {
	if !key.HasTradePoint {
		return nil
	}
	return key.TradePointID
}
