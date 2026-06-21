package writefirms

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

func (s *Service) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	firm, err := s.repo.UpdateFirm(ctx, id, name)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.Firm{}, fmt.Errorf("%w: фирма с id %d не найдена", models.ErrNotFound, id)
		}
		return quik.Firm{}, err
	}
	return firm, nil
}
