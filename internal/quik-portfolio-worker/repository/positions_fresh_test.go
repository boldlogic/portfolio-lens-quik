package repository

import (
	"strings"
	"testing"
)

func TestPositionSQLKeepsOtherSource(t *testing.T) {
	t.Parallel()
	for _, query := range []string{deletePositionsWithoutLimitsSQL, replacePositionsForDateSQL} {
		if !strings.Contains(query, "ext_system_id IS NULL OR p.ext_system_id = @quik") {
			t.Fatal("удаление не ограничено источником QUIK")
		}
	}
	if !strings.Contains(replacePositionsForDateSQL, "p.source_updated_at < @observed") {
		t.Fatal("обновление не смотрит время")
	}
	if !strings.Contains(replacePositionsForDateSQL, "EXISTS (SELECT p.balance EXCEPT SELECT s.balance)") {
		t.Fatal("одинаковый остаток не должен обновлять время")
	}
}
