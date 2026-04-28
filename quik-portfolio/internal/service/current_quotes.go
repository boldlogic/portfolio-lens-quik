package service

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
)

// нужно потом вынести в отдельный сервис когда появится время
func (s *Service) GetCurrentQuotes(ctx context.Context) ([]models.CurrentQuote, error) {
	return s.repo.SelectCurrentQuotes(ctx)
}

func (s *Service) GetCurrentQuotesForKeys(ctx context.Context, keys []string) ([]models.CurrentQuote, error) {
	return s.repo.SelectCurrentQuotesForKeys(ctx, keys)
}
