// Package recordstore persists httpstream records in SQLite.
package recordstore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  payload TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_records_id ON records(id);
`

// Store provides SQLite-backed record persistence.
type Store struct {
	db *sql.DB
}

// Open opens or creates a record database at path.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("init schema: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database.
func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// Insert stores a record payload and returns its row id.
func (s *Store) Insert(payload []byte) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO records (payload) VALUES (?)`, string(payload))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetRecent returns the most recent limit records as decoded JSON values.
func (s *Store) GetRecent(limit int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := s.db.Query(`
		SELECT payload FROM records
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payloads []string
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}
		payloads = append(payloads, payload)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse to chronological order (oldest first among the batch).
	for i, j := 0, len(payloads)-1; i < j; i, j = i+1, j-1 {
		payloads[i], payloads[j] = payloads[j], payloads[i]
	}

	return decodePayloads(payloads)
}

// GetSince returns records with id > afterID, up to limit, in ascending id order.
func (s *Store) GetSince(afterID int64, limit int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 1000
	}

	rows, err := s.db.Query(`
		SELECT payload FROM records
		WHERE id > ?
		ORDER BY id ASC
		LIMIT ?
	`, afterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payloads []string
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}
		payloads = append(payloads, payload)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return decodePayloads(payloads)
}

// LatestID returns the highest record id, or 0 if empty.
func (s *Store) LatestID() (int64, error) {
	var id sql.NullInt64
	err := s.db.QueryRow(`SELECT MAX(id) FROM records`).Scan(&id)
	if err != nil {
		return 0, err
	}
	if !id.Valid {
		return 0, nil
	}
	return id.Int64, nil
}

// Count returns total record count.
func (s *Store) Count() (int64, error) {
	var n int64
	err := s.db.QueryRow(`SELECT COUNT(*) FROM records`).Scan(&n)
	return n, err
}

func decodePayloads(payloads []string) ([]interface{}, error) {
	out := make([]interface{}, 0, len(payloads))
	for _, p := range payloads {
		var v interface{}
		if err := json.Unmarshal([]byte(p), &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}
