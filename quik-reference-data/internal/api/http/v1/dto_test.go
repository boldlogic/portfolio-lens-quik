package v1

import (
	"fmt"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirmToResp(t *testing.T) {
	got := firmToResp(quik.Firm{Id: 1, Code: "MC0000000000", Name: "Брокер"})

	assert.Equal(t, firmRespDTO{
		Id:   1,
		Code: "MC0000000000",
		Name: "Брокер",
	}, got, "ожидали полное соответствие")
}

func TestFirmsToResp(t *testing.T) {
	t.Run("nil_слайс_возвращает_пустой_ответ", func(t *testing.T) {
		got := firmsToResp(nil)
		fmt.Println(got)

		require.NotNil(t, got, "ожидали не nil")
		assert.Empty(t, got, "ожидали пустой слайс")
		assert.Equal(t, []firmRespDTO{}, got, "ожидали пустой слайс")
	})

	t.Run("пустой_слайс_возвращает_пустой_ответ", func(t *testing.T) {
		got := firmsToResp([]quik.Firm{})
		fmt.Println(got)

		require.NotNil(t, got, "ожидали не nil")
		assert.Empty(t, got, "ожидали пустой слайс")
		assert.Equal(t, []firmRespDTO{}, got, "ожидали пустой слайс")
	})

	t.Run("несколько_фирм", func(t *testing.T) {
		got := firmsToResp([]quik.Firm{
			{Id: 1, Code: "MC0000000000", Name: "Первый брокер"},
			{Id: 2, Code: "SP0000000000", Name: "Второй брокер"},
		})
		assert.Len(t, got, 2)

		assert.Equal(t, []firmRespDTO{
			{Id: 1, Code: "MC0000000000", Name: "Первый брокер"},
			{Id: 2, Code: "SP0000000000", Name: "Второй брокер"},
		}, got)
	})
}
