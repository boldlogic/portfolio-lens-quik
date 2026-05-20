package repository

import (
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

func (r *Repository) logWrapper(f string, date time.Time, err error) {

	if err == nil {
		r.Logger.Debug("успех", zap.String("func", f), zap.Time("load_date", date))
		return
	}

	if shutdown.IsExceeded(err) {
		return
	}
	r.Logger.Error("ошибка выполнения", zap.String("func", f), zap.Time("load_date", date), zap.Error(err))

}

func (r *Repository) finalizeSelectErr(funcName string, date time.Time, err error) error {
	r.logWrapper(funcName, date, err)
	if err != nil && !shutdown.IsExceeded(err) {
		return models.ErrRetrievingData
	}
	return err
}
