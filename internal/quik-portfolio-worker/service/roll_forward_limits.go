package service

import "context"

func (s *Service) DoRollForwardMoneyLimits(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectMoneyLimitsMaxDate,
		s.repo.InsertMoneyLimitsCopy,
		s.repo.DeleteMoneyLimits,
	)
}

func (s *Service) DoRollForwardSecurityLimits(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectSecurityLimitsMaxDate,
		s.repo.InsertSecurityLimitsCopy,
		s.repo.DeleteSecurityLimits,
	)
}

func (s *Service) DoRollForwardOtc(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectSecurityLimitsOtcMaxDate,
		s.repo.InsertSecurityLimitsOtcCopy,
		s.repo.DeleteSecurityLimitsOtc,
	)
}
