package service

import "context"

func (s *Service) RebuildPositions(ctx context.Context) error {
	dates, err := s.repo.ListPositionLoadDates(ctx)
	if err != nil {
		return err
	}
	for _, date := range dates {
		if err := s.repo.ReplacePositionsForDate(ctx, date); err != nil {
			return err
		}
	}
	return s.repo.DeletePositionsWithoutLimits(ctx)
}
