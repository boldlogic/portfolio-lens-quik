package service

import (
	"errors"
	"testing"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/quik-portfolio/pkg/models"
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

func TestMinRollForwardDate_nil_совпадает_с_Today(t *testing.T) {
	got := minRollForwardDate(nil)
	want := dates.Today()
	if dates.DateToYYYYMMDD(got) != dates.DateToYYYYMMDD(want) {
		t.Fatalf("ожидали сегодняшнюю дату: got=%v want=%v", got, want)
	}
}
