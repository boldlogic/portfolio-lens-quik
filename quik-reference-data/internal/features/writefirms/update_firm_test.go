package writefirms

import (
	"context"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUpdateFirm(t *testing.T) {
	t.Run("успешное_обновление", func(t *testing.T) {
		svc := NewService(repo{}, zap.NewNop())
		got, err := svc.UpdateFirm(context.Background(), 1, "ФирмаНовая")
		require.Equal(t, quik.Firm{Id: 1, Code: "code", Name: "ФирмаНовая"}, got)
		require.NoError(t, err)
	})
	t.Run("ErrNotFound", func(t *testing.T) {
		svc := NewService(repo{err: models.ErrNotFound}, zap.NewNop())
		got, err := svc.UpdateFirm(context.Background(), 1, "ФирмаНовая")
		require.Equal(t, quik.Firm{}, got)
		require.ErrorIs(t, err, models.ErrNotFound)
	})
}
