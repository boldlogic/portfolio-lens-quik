package sqlserver

import (
	"context"
	"fmt"
)

const linkQuoteSQL = `
UPDATE quik.current_quotes
SET instrument_id = @p1
WHERE CONVERT(binary(8), rw) = @p2
	AND (instrument_id IS NULL OR instrument_id <> @p1);`

func (r *Repository) LinkQuote(ctx context.Context, sourceRowVersion []byte, instrumentID int64) (bool, error) {
	res, err := r.Db.ExecContext(ctx, linkQuoteSQL, instrumentID, sourceRowVersion)
	if err != nil {
		return false, fmt.Errorf("ошибка привязки котировки к инструменту: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("ошибка проверки привязки котировки к инструменту: %w", err)
	}
	return affected > 0, nil
}
