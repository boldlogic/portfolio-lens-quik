package writefirms

import (
	"context"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type repo struct {
	err error
}

func (r repo) InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	return quik.Firm{
		Id:   1,
		Code: code,
		Name: name,
	}, r.err
}

func (r repo) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	return quik.Firm{
		Id:   id,
		Code: "code",
		Name: name,
	}, r.err
}

func TestCreateFirm(t *testing.T) {
	t.Run("успешное_создание", func(t *testing.T) {
		svc := NewService(repo{}, zap.NewNop())
		got, err := svc.CreateFirm(context.Background(), "firm", "Фирма")
		require.Equal(t, quik.Firm{Id: 1, Code: "firm", Name: "Фирма"}, got)
		require.NoError(t, err)
	})
	t.Run("конфликт", func(t *testing.T) {
		svc := NewService(repo{err: models.ErrConflict}, zap.NewNop())
		got, err := svc.CreateFirm(context.Background(), "firm", "Фирма")
		require.Equal(t, quik.Firm{}, got)
		require.ErrorIs(t, err, models.ErrConflict)
	})
}
