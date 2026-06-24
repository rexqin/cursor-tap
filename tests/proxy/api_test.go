package proxy_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	appconfig "github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/proxy"
	"github.com/stretchr/testify/require"
)

func newTestProxyServer(t *testing.T, withRecorder bool) *proxy.Server {
	t.Helper()

	dir := t.TempDir()
	cfg := appconfig.Config{
		HTTPPort:   18080,
		SOCKS5Port: 11080,
		APIPort:    19090,
		CertDir:    dir,
		DataDir:    filepath.Join(dir, "data"),
	}
	if withRecorder {
		cfg.EnableHTTPParsing = true
		cfg.HTTPRecordFile = filepath.Join(dir, "access.log")
	}

	server, err := proxy.NewServer(cfg)
	require.NoError(t, err)
	t.Cleanup(server.Close)
	return server
}

func testAPIServer(t *testing.T, server *proxy.Server) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	server.RegisterAPIRoutes(mux)
	return httptest.NewServer(mux)
}

func TestRegisterAPIRoutes_status(t *testing.T) {
	server := newTestProxyServer(t, false)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/status")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))

	var body map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "running", body["status"])
}

func TestRegisterAPIRoutes_stats(t *testing.T) {
	server := newTestProxyServer(t, false)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/stats")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.EqualValues(t, 0, body["active_sessions"])
	require.EqualValues(t, 0, body["total_sessions"])
	require.EqualValues(t, 0, body["ws_clients"])
}

func TestRegisterAPIRoutes_caCert(t *testing.T) {
	server := newTestProxyServer(t, false)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/ca/cert")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "BEGIN CERTIFICATE")
}

func TestRegisterAPIRoutes_sessionsNotImplemented(t *testing.T) {
	server := newTestProxyServer(t, false)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/sessions")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRegisterAPIRoutes_recordsWithoutRecorder(t *testing.T) {
	server := newTestProxyServer(t, false)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/records")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRegisterAPIRoutes_recordsWithRecorder(t *testing.T) {
	server := newTestProxyServer(t, true)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/records?limit=100")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var records []interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&records))
	require.NotNil(t, records)
}

func TestRegisterAPIRoutes_wsRecordsWithoutRecorder(t *testing.T) {
	server := newTestProxyServer(t, false)
	ts := testAPIServer(t, server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ws/records")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
