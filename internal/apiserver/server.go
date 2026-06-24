// Package apiserver provides the standalone management API HTTP server.
package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/burpheart/cursor-tap/internal/api"
	appconfig "github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/notify"
	"github.com/burpheart/cursor-tap/internal/recordstore"
)

// Server is the standalone API server process.
type Server struct {
	config   appconfig.APIConfig
	store    *recordstore.Store
	hub      *api.Hub
	certPath string

	httpServer *http.Server
	lastID     atomic.Int64

	mu      sync.Mutex
	running bool
}

// New creates a new API server.
func New(config appconfig.APIConfig) (*Server, error) {
	if err := os.MkdirAll(config.CertDir, 0755); err != nil {
		return nil, fmt.Errorf("create cert dir: %w", err)
	}

	store, err := recordstore.Open(config.RecordDB)
	if err != nil {
		return nil, fmt.Errorf("open record store: %w", err)
	}

	certPath := filepath.Join(config.CertDir, "ca", "ca.crt")

	return &Server{
		config:   config,
		store:    store,
		hub:      api.NewHub(),
		certPath: certPath,
	}, nil
}

// Start starts the API HTTP server (blocking).
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("server already running")
	}
	s.running = true
	s.mu.Unlock()

	addr := fmt.Sprintf("127.0.0.1:%d", s.config.Port)
	addrFile := filepath.Join(s.config.CertDir, "api.addr")
	if err := os.WriteFile(addrFile, []byte(addr), 0644); err != nil {
		fmt.Printf("[WARN] Failed to write API address file: %v\n", err)
	}

	go s.hub.Run()

	if latest, err := s.store.LatestID(); err == nil {
		s.lastID.Store(latest)
	}

	mux := http.NewServeMux()
	s.registerRoutes(mux)

	s.httpServer = &http.Server{Addr: addr, Handler: mux}
	fmt.Printf("[INFO] API server listening on %s\n", addr)
	fmt.Printf("[INFO] Record DB: %s\n", s.config.RecordDB)

	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Stop shuts down the API server.
func (s *Server) Stop() error {
	s.mu.Lock()
	wasRunning := s.running
	s.running = false
	s.mu.Unlock()

	if !wasRunning {
		return s.Close()
	}

	var err error
	if s.httpServer != nil {
		err = s.httpServer.Close()
	}

	addrFile := filepath.Join(s.config.CertDir, "api.addr")
	os.Remove(addrFile)

	if closeErr := s.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	return err
}

// Close releases resources.
func (s *Server) Close() error {
	if s.store != nil {
		err := s.store.Close()
		s.store = nil
		return err
	}
	return nil
}

// RegisterRoutes registers API routes on mux (for testing).
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	s.registerRoutes(mux)
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(`{"status":"running"}`))
	})

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		count, _ := s.store.Count()
		stats := fmt.Sprintf(`{"active_sessions":0,"total_sessions":0,"total_bytes_sent":0,"total_bytes_received":0,"ws_clients":%d,"record_count":%d}`,
			s.hub.ClientCount(), count)
		w.Write([]byte(stats))
	})

	mux.HandleFunc("/api/ca/cert", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.certPath)
	})

	mux.HandleFunc("/internal/notify", s.handleNotify)

	handler := api.NewHandler(s.hub, &storeAdapter{store: s.store})
	handler.RegisterRoutes(mux)
}

func (s *Server) handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || (host != "127.0.0.1" && host != "::1" && host != "[::1]") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var payload notify.Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	after := s.lastID.Load()
	if payload.LatestID > after {
		records, err := s.store.GetSince(after, 1000)
		if err == nil {
			for _, rec := range records {
				s.hub.Broadcast(rec)
			}
		}
		s.lastID.Store(payload.LatestID)
	}

	w.WriteHeader(http.StatusNoContent)
}

type storeAdapter struct {
	store *recordstore.Store
}

func (a *storeAdapter) GetRecentRecords(limit int) []interface{} {
	records, err := a.store.GetRecent(limit)
	if err != nil {
		return []interface{}{}
	}
	return records
}
