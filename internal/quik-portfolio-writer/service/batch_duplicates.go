package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func duplicateLineGroupsFromHash(byHash map[[32]byte][]uint) [][]uint {
	groups := make([][]uint, 0)
	for _, lines := range byHash {
		if len(lines) < 2 {
			continue
		}
		sort.Slice(lines, func(i, j int) bool { return lines[i] < lines[j] })
		groups = append(groups, lines)
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	return groups
}

func formatDuplicateLimitLines(groups [][]uint) string {
	parts := make([]string, 0, len(groups))
	for _, lines := range groups {
		strLines := make([]string, 0, len(lines))
		for _, line := range lines {
			strLines = append(strLines, strconv.FormatUint(uint64(line), 10))
		}
		parts = append(parts, "строки "+strings.Join(strLines, ", "))
	}
	return strings.Join(parts, "; ")
}

func errDuplicateLimits(groups [][]uint) error {
	return fmt.Errorf("%w: дубликаты лимитов: %s", models.ErrBusinessValidation, formatDuplicateLimitLines(groups))
}
