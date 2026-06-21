package readfirms

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

func (s *service) GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	firm, err := s.repo.SelectFirmByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.Firm{}, fmt.Errorf("%w: фирма с id %d не найдена", models.ErrNotFound, id)
		}
		return quik.Firm{}, err
	}
	return firm, nil
}
