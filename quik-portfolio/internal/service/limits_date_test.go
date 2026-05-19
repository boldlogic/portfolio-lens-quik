package service

import (
	"errors"
	"testing"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func TestCheckLimitDate_после_сегодня_ErrBusinessValidation(t *testing.T) {
	today := dates.Today()
	min := today.AddDate(-1, 0, 0)
	future := today.AddDate(0, 0, 1)
	err := checkLimitDate(future, min)
	if !errors.Is(err, models.ErrBusinessValidation) {
		t.Fatalf("ожидали ErrBusinessValidation: %v", err)
	}
}
