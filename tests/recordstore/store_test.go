package recordstore_test

import (
	"encoding/json"
	"path/filepath"
	"sync"
	"testing"

	"github.com/burpheart/cursor-tap/internal/recordstore"
	"github.com/stretchr/testify/require"
)

func TestStoreInsertAndGetRecent(t *testing.T) {
	dir := t.TempDir()
	store, err := recordstore.Open(filepath.Join(dir, "records.db"))
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })

	for i := 1; i <= 3; i++ {
		payload, err := json.Marshal(map[string]any{"type": "request", "n": i})
		require.NoError(t, err)
		id, err := store.Insert(payload)
		require.NoError(t, err)
		require.Equal(t, int64(i), id)
	}

	recent, err := store.GetRecent(2)
	require.NoError(t, err)
	require.Len(t, recent, 2)

	first := recent[0].(map[string]any)
	second := recent[1].(map[string]any)
	require.EqualValues(t, 2, first["n"])
	require.EqualValues(t, 3, second["n"])
}

func TestStoreGetSince(t *testing.T) {
	dir := t.TempDir()
	store, err := recordstore.Open(filepath.Join(dir, "records.db"))
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })

	for i := 0; i < 5; i++ {
		payload, _ := json.Marshal(map[string]any{"i": i})
		_, err := store.Insert(payload)
		require.NoError(t, err)
	}

	since, err := store.GetSince(2, 10)
	require.NoError(t, err)
	require.Len(t, since, 3)
}

func TestStoreConcurrentInsert(t *testing.T) {
	dir := t.TempDir()
	store, err := recordstore.Open(filepath.Join(dir, "records.db"))
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			payload, _ := json.Marshal(map[string]any{"n": n})
			_, err := store.Insert(payload)
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	count, err := store.Count()
	require.NoError(t, err)
	require.Equal(t, int64(20), count)
}
