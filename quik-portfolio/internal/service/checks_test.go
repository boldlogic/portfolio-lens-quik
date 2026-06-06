package service

import (
	"testing"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_dedublicateClientCodes(t *testing.T) {
	tests := []struct {
		name        string
		clientCodes []string
		want        []string
		wantErr     bool
	}{
		{
			name:        "один_клиент",
			clientCodes: []string{"AA"},
			want:        []string{"AA"},
			wantErr:     false,
		},
		{
			name:        "10_клиентов",
			clientCodes: []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10"},
			want:        []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10"},
			wantErr:     false,
		},
		{
			name:        "регистр",
			clientCodes: []string{"aa"},
			want:        []string{"AA"},
			wantErr:     false,
		},
		{
			name:        "nil",
			clientCodes: nil,
			want:        nil,
			wantErr:     false,
		},
		{
			name:        "пустой",
			clientCodes: []string{},
			want:        nil,
			wantErr:     false,
		},
		{
			name:        "пробелы",
			clientCodes: []string{" aa "},
			want:        []string{"AA"},
			wantErr:     false,
		},
		{
			name:        "дубликаты",
			clientCodes: []string{"a1", "A1"},
			want:        []string{"A1"},
			wantErr:     false,
		},
		{
			name:        "10_уникальных",
			clientCodes: []string{"A1", "a1", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10"},
			want:        []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10"},
			wantErr:     false,
		},
		{
			name:        "слишком_много",
			clientCodes: []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11"},
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "длинная_строка",
			clientCodes: []string{"1234567890123"},
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			got, gotErr := dedublicateClientCodes(tt.clientCodes)
			if tt.wantErr {
				require.ErrorIs(t, gotErr, models.ErrBusinessValidation)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_checkLimitDate(t *testing.T) {
	today := dates.Today()

	tests := []struct {
		name    string
		date    time.Time
		wantErr bool
	}{
		{
			name:    "вчера",
			date:    today.AddDate(0, 0, -1),
			wantErr: false,
		},
		{
			name:    "сегодня",
			date:    today,
			wantErr: false,
		},
		{
			name:    "завтра",
			date:    today.AddDate(0, 0, 1),
			wantErr: false,
		},
		{
			name:    "слишком_рано",
			date:    today.AddDate(0, 0, -2),
			wantErr: true,
		},
		{
			name:    "слишком_поздно",
			date:    today.AddDate(0, 0, 2),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := checkLimitDate(tt.date)
			if tt.wantErr {
				require.ErrorIs(t, gotErr, models.ErrBusinessValidation)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}

func Test_validateLimitsContract(t *testing.T) {
	today := dates.Today()

	tests := []struct {
		name        string
		date        int
		clientCodes []string
		want        []string
		wantErr     bool
	}{
		{
			name:        "валидный_контракт",
			date:        0,
			clientCodes: []string{" a1 ", "A1", "b2"},
			want:        []string{"A1", "B2"},
			wantErr:     false,
		},
		{
			name:        "невалидная_дата",
			date:        -2,
			clientCodes: []string{"A1"},
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "невалидный_код_клиента",
			date:        0,
			clientCodes: []string{"1234567890123"},
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := validateLimitsContract(today.AddDate(0, 0, tt.date), tt.clientCodes)
			if tt.wantErr {
				require.ErrorIs(t, gotErr, models.ErrBusinessValidation)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
