package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/burpheart/cursor-tap/internal/api"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type mockRecordStore struct {
	limit   int
	records []interface{}
}

func (m *mockRecordStore) GetRecentRecords(limit int) []interface{} {
	m.limit = limit
	if m.records == nil {
		return []interface{}{}
	}
	return m.records
}

func newTestHandler(t *testing.T, store api.RecordStore) (*api.Handler, *api.Hub) {
	t.Helper()
	hub := api.NewHub()
	go hub.Run()
	t.Cleanup(func() {
		// hub has no Stop; goroutine exits with test process
	})
	return api.NewHandler(hub, store), hub
}

func TestHandleGetRecords_defaultLimit(t *testing.T) {
	store := &mockRecordStore{
		records: []interface{}{
			map[string]string{"type": "request"},
		},
	}
	handler, _ := newTestHandler(t, store)

	req := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	rec := httptest.NewRecorder()
	handler.HandleGetRecords(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, 100, store.limit)

	var got []map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Len(t, got, 1)
	require.Equal(t, "request", got[0]["type"])
}

func TestHandleGetRecords_customLimit(t *testing.T) {
	store := &mockRecordStore{}
	handler, _ := newTestHandler(t, store)

	req := httptest.NewRequest(http.MethodGet, "/api/records?limit=50", nil)
	rec := httptest.NewRecorder()
	handler.HandleGetRecords(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, 50, store.limit)
}

func TestHandleGetRecords_limitValidation(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		wantLimit int
	}{
		{"invalid uses default", "limit=abc", 100},
		{"zero uses default", "limit=0", 100},
		{"negative uses default", "limit=-1", 100},
		{"over max uses default", "limit=5000", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockRecordStore{}
			handler, _ := newTestHandler(t, store)

			req := httptest.NewRequest(http.MethodGet, "/api/records?"+tt.query, nil)
			rec := httptest.NewRecorder()
			handler.HandleGetRecords(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, tt.wantLimit, store.limit)
		})
	}
}

func TestHandleCORS(t *testing.T) {
	handler, _ := newTestHandler(t, &mockRecordStore{})

	req := httptest.NewRequest(http.MethodOptions, "/api/records", nil)
	rec := httptest.NewRecorder()
	handler.HandleCORS(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type", rec.Header().Get("Access-Control-Allow-Headers"))
}

func TestRegisterRoutes_optionsPreflight(t *testing.T) {
	store := &mockRecordStore{}
	handler, _ := newTestHandler(t, store)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	req, err := http.NewRequest(http.MethodOptions, server.URL+"/api/records", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestRegisterRoutes_getRecords(t *testing.T) {
	store := &mockRecordStore{
		records: []interface{}{
			map[string]string{"type": "grpc"},
		},
	}
	handler, _ := newTestHandler(t, store)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/api/records?limit=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 10, store.limit)

	var got []map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Len(t, got, 1)
	require.Equal(t, "grpc", got[0]["type"])
}

func TestHandleWebSocket_upgrade(t *testing.T) {
	handler, _ := newTestHandler(t, &mockRecordStore{})

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/records"
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
	conn.Close()
}
