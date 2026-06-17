package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func jsonObjectKeys(t *testing.T, v any) map[string]json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	var out map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(b, &out))
	return out
}
