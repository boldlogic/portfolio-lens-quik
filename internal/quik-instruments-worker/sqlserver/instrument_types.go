package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

const selectInstrumentTypeIDSQL = `
SELECT type_id
FROM ref.instrument_types
WHERE title = @p1;`

func (r *Repository) SelectInstrumentTypeID(ctx context.Context, title quik.InstrumentTypeTitle) (uint8, error) {
	var typeID uint8
	err := r.Db.QueryRowContext(ctx, selectInstrumentTypeIDSQL, string(title)).Scan(&typeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("тип инструмента не найден: %s", title)
		}
		return 0, fmt.Errorf("ошибка чтения типа инструмента: %w", err)
	}
	return typeID, nil
}
