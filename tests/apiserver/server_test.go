package apiserver_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/burpheart/cursor-tap/internal/apiserver"
	appconfig "github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/notify"
	"github.com/burpheart/cursor-tap/internal/recordstore"
	"github.com/stretchr/testify/require"
)

func setupAPIServer(t *testing.T) (*apiserver.Server, string, string) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "records.db")

	caDir := filepath.Join(dir, "ca")
	require.NoError(t, os.MkdirAll(caDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(caDir, "ca.crt"), []byte("-----BEGIN CERTIFICATE-----\nTEST\n-----END CERTIFICATE-----\n"), 0644))

	cfg := appconfig.APIConfig{
		Port:     9090,
		CertDir:  dir,
		RecordDB: dbPath,
	}

	server, err := apiserver.New(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { server.Close() })

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	return server, dbPath, ts.URL
}

func insertRecord(t *testing.T, dbPath string, rec map[string]string) int64 {
	t.Helper()
	store, err := recordstore.Open(dbPath)
	require.NoError(t, err)
	defer store.Close()

	payload, err := json.Marshal(rec)
	require.NoError(t, err)
	id, err := store.Insert(payload)
	require.NoError(t, err)
	return id
}

func TestRegisterAPIRoutes_status(t *testing.T) {
	_, _, baseURL := setupAPIServer(t)

	resp, err := http.Get(baseURL + "/api/status")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "running", body["status"])
}

func TestRegisterAPIRoutes_stats(t *testing.T) {
	_, dbPath, baseURL := setupAPIServer(t)

	insertRecord(t, dbPath, map[string]string{"type": "grpc"})

	resp, err := http.Get(baseURL + "/api/stats")
	require.NoError(t, err)
	defer resp.Body.Close()

	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.EqualValues(t, 1, body["record_count"])
}

func TestRegisterAPIRoutes_caCert(t *testing.T) {
	_, _, baseURL := setupAPIServer(t)

	resp, err := http.Get(baseURL + "/api/ca/cert")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "BEGIN CERTIFICATE")
}

func TestRegisterAPIRoutes_records(t *testing.T) {
	_, dbPath, baseURL := setupAPIServer(t)

	insertRecord(t, dbPath, map[string]string{"type": "request"})

	resp, err := http.Get(baseURL + "/api/records?limit=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var records []map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&records))
	require.Len(t, records, 1)
	require.Equal(t, "request", records[0]["type"])
}

func TestRegisterAPIRoutes_sessionsNotImplemented(t *testing.T) {
	_, _, baseURL := setupAPIServer(t)

	resp, err := http.Get(baseURL + "/api/sessions")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestInternalNotify_broadcastsRecords(t *testing.T) {
	_, dbPath, baseURL := setupAPIServer(t)

	id := insertRecord(t, dbPath, map[string]string{"type": "grpc"})

	body, _ := json.Marshal(notify.Payload{LatestID: id})
	resp, err := http.Post(baseURL+"/internal/notify", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp, err = http.Get(baseURL + "/api/records?limit=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	var records []map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&records))
	require.Len(t, records, 1)
}
