package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/burpheart/cursor-tap/internal/notify"
	"github.com/stretchr/testify/require"
)

func TestClientNotifyLatest(t *testing.T) {
	var received atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var payload notify.Payload
		require.NoError(t, json.Unmarshal(body, &payload))
		require.Equal(t, int64(42), payload.LatestID)
		received.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client := notify.NewClient(server.URL)
	client.NotifyLatest(42)

	require.Eventually(t, func() bool {
		return received.Load() == 1
	}, time.Second, 10*time.Millisecond)
}

func TestClientEmptyURLNoOp(t *testing.T) {
	client := notify.NewClient("")
	require.NotPanics(t, func() {
		client.NotifyLatest(1)
	})
}
