package writefirms

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	firm, err := s.repo.InsertFirm(ctx, code, name)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return quik.Firm{}, fmt.Errorf("%w: фирма c firmCode %s уже существует", models.ErrConflict, code)
		}
		return quik.Firm{}, err
	}
	return firm, nil
}
