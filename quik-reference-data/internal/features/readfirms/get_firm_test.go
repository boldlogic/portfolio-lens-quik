package readfirms

import (
	"context"
	"errors"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type repoStub struct {
	err error
}

func (r repoStub) SelectFirms(ctx context.Context) ([]quik.Firm, error) {
	return nil, nil
}

func (r repoStub) SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	if r.err != nil {
		return quik.Firm{}, r.err
	}
	return quik.Firm{Id: id, Code: "firm", Name: "Фирма"}, nil
}

func TestGetFirmByID_фирма_не_найдена(t *testing.T) {
	svc := NewService(repoStub{err: models.ErrNotFound}, zap.NewNop())
	_, err := svc.GetFirmByID(context.Background(), 1)
	if !errors.Is(err, models.ErrNotFound) {
		t.Fatalf("ожидали ErrNotFound: %v", err)
	}
}
