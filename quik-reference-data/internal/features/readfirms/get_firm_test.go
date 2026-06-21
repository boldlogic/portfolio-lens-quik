package readfirms

import (
	"context"
	"errors"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type repo struct {
	err error
}

func (r repo) SelectFirms(ctx context.Context) ([]quik.Firm, error) {
	return nil, nil
}

func (r repo) SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	return quik.Firm{Id: id, Code: "firm", Name: "Фирма"}, r.err
}

func TestGetFirmByID(t *testing.T) {
	t.Run("успешно_возвращает_фирму", func(t *testing.T) {
		svc := NewService(repo{}, zap.NewNop())
		got, err := svc.GetFirmByID(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, quik.Firm{Id: 1, Code: "firm", Name: "Фирма"}, got)
	})
	t.Run("фирма_не_найдена", func(t *testing.T) {
		svc := NewService(repo{err: models.ErrNotFound}, zap.NewNop())
		got, err := svc.GetFirmByID(context.Background(), 1)
		require.ErrorIs(t, err, models.ErrNotFound)
		require.ErrorContains(t, err, "фирма с id 1 не найдена")
		require.Equal(t, quik.Firm{}, got)
	})

	t.Run("неизвестная_ошибка", func(t *testing.T) {
		unexpectedErr := errors.New("db error")
		svc := NewService(repo{err: unexpectedErr}, zap.NewNop())
		got, err := svc.GetFirmByID(context.Background(), 1)
		require.ErrorIs(t, err, unexpectedErr)
		require.Equal(t, quik.Firm{}, got)
	})

}
