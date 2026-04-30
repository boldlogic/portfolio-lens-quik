package readfirms

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *service) GetFirms(ctx context.Context) ([]quik.Firm, error) {
	return s.repo.SelectFirms(ctx)
}
