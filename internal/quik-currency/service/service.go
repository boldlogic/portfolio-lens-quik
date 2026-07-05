package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/packages/shutdown"
	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	repo           Repo
	currencyWorker *currencyWorker
}

type Repo interface {
	MergeFxCBRRates(ctx context.Context) error
	currencyRepo
}

func (s *Service) Run(ctx context.Context) error {
	return s.currencyWorker.run(ctx)
}

func NewService(logger *zap.Logger, repo Repo, currencyCfg CurrencyWorkerConfig) (*Service, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	if repo == nil {
		return nil, fmt.Errorf("не передан репозиторий")
	}
	return &Service{
		logger:         logger,
		repo:           repo,
		currencyWorker: newCurrencyWorker(currencyCfg, repo, logger),
	}, nil
}

func (s *Service) logError(msg string, err error) {
	if shutdown.IsExceeded(err) {
		return
	}
	if err != nil {
		s.logger.Error(msg, zap.Error(err))
	}
}
