package httpstream

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/burpheart/cursor-tap/internal/recordstore"
)

// Record represents a persisted traffic record.
type Record struct {
	Timestamp   string `json:"ts"`
	SessionID   string `json:"session"`
	SessionSeq  int64  `json:"seq"`
	RecordIndex int64  `json:"index"`
	Type        string `json:"type"`

	Method string `json:"method,omitempty"`
	URL    string `json:"url,omitempty"`
	Host   string `json:"host,omitempty"`

	Status     int    `json:"status,omitempty"`
	StatusText string `json:"status_text,omitempty"`

	EventType string `json:"event_type,omitempty"`
	EventID   string `json:"event_id,omitempty"`
	EventData string `json:"event_data,omitempty"`

	Headers map[string][]string `json:"headers,omitempty"`

	Direction    string `json:"direction,omitempty"`
	Size         int    `json:"size,omitempty"`
	Body         string `json:"body,omitempty"`
	BodyBase64   string `json:"body_base64,omitempty"`
	BodyEncoding string `json:"body_encoding,omitempty"`
	ContentType  string `json:"content_type,omitempty"`

	GRPCService    string `json:"grpc_service,omitempty"`
	GRPCMethod     string `json:"grpc_method,omitempty"`
	GRPCData       string `json:"grpc_data,omitempty"`
	GRPCStreaming  bool   `json:"grpc_streaming,omitempty"`
	GRPCFrameIndex int    `json:"grpc_frame_index,omitempty"`
	GRPCCompressed bool   `json:"grpc_compressed,omitempty"`
	GRPCRawData    string `json:"grpc_raw,omitempty"`

	Error string `json:"error,omitempty"`
}

// RecordCallback is called after a record is persisted.
type RecordCallback func(id int64, rec Record)

// Recorder writes HTTP traffic to SQLite with session tracking.
type Recorder struct {
	store    *recordstore.Store
	logLevel LogLevel

	records    atomic.Int64
	sessionSeq atomic.Int64
	onRecord   RecordCallback
}

// RecorderOption configures a Recorder.
type RecorderOption func(*Recorder)

// WithRecorderLogLevel sets the log level for recording.
func WithRecorderLogLevel(level LogLevel) RecorderOption {
	return func(r *Recorder) { r.logLevel = level }
}

// WithOnRecord sets a callback for each record written.
func WithOnRecord(cb RecordCallback) RecorderOption {
	return func(r *Recorder) { r.onRecord = cb }
}

// NewRecorder creates a SQLite-backed recorder at dbPath.
func NewRecorder(dbPath string, opts ...RecorderOption) (*Recorder, error) {
	store, err := recordstore.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open record store: %w", err)
	}

	r := &Recorder{
		store:    store,
		logLevel: LogLevelBasic,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// Close closes the underlying store.
func (r *Recorder) Close() error {
	if r.store == nil {
		return nil
	}
	return r.store.Close()
}

// Store returns the underlying record store (for tests).
func (r *Recorder) Store() *recordstore.Store {
	return r.store
}

func (r *Recorder) write(rec Record) error {
	payload, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	id, err := r.store.Insert(payload)
	if err != nil {
		return err
	}

	r.records.Add(1)

	if r.onRecord != nil {
		r.onRecord(id, rec)
	}

	return nil
}

// RecordCount returns the number of records written.
func (r *Recorder) RecordCount() int64 {
	return r.records.Load()
}

// Session represents a tracked HTTP session.
type Session struct {
	ID          string
	Seq         int64
	Host        string
	recorder    *Recorder
	recordIndex int64
}

// NewSession creates a new tracked session.
func (r *Recorder) NewSession(host string) *Session {
	seq := r.sessionSeq.Add(1)
	return &Session{
		ID:          generateSessionID(),
		Seq:         seq,
		Host:        host,
		recorder:    r,
		recordIndex: 0,
	}
}

func (s *Session) nextRecordIndex() int64 {
	s.recordIndex++
	return s.recordIndex
}

func timestamp() string {
	return time.Now().Format(time.RFC3339Nano)
}

func (s *Session) LogRequest(msg *HTTPMessage) {
	req := msg.Request
	if req == nil {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "request",
		Method:      req.Method,
		URL:         req.URL.RequestURI(),
		Host:        s.Host,
		Headers:     cloneHeaders(req.Header),
		ContentType: req.Header.Get("Content-Type"),
	}

	s.recorder.write(rec)
}

func (s *Session) LogResponse(msg *HTTPMessage) {
	resp := msg.Response
	if resp == nil {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "response",
		Status:      resp.StatusCode,
		StatusText:  resp.Status,
		Host:        s.Host,
		Headers:     cloneHeaders(resp.Header),
		ContentType: resp.Header.Get("Content-Type"),
	}

	s.recorder.write(rec)
}

func (s *Session) LogSSE(host string, event *SSEEvent) {
	eventType := event.Event
	if eventType == "" {
		eventType = "message"
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "sse",
		Host:        host,
		EventType:   eventType,
		EventID:     event.ID,
		EventData:   truncateString(event.Data, 1000),
	}

	s.recorder.write(rec)
}

func (s *Session) LogBody(dir Direction, host string, data []byte) {
	if len(data) == 0 {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "body",
		Direction:   dir.String(),
		Host:        host,
		Size:        len(data),
	}

	if utf8.Valid(data) && isPrintableText(data) {
		rec.Body = string(data)
		rec.BodyEncoding = "text"
	} else {
		rec.BodyBase64 = base64.StdEncoding.EncodeToString(data)
		rec.BodyEncoding = "base64"
	}

	s.recorder.write(rec)
}

func (s *Session) LogGRPC(msg *GRPCMessage) {
	rec := Record{
		Timestamp:      timestamp(),
		SessionID:      s.ID,
		SessionSeq:     s.Seq,
		RecordIndex:    s.nextRecordIndex(),
		Type:           "grpc",
		Direction:      msg.Direction.String(),
		Host:           s.Host,
		GRPCService:    msg.Service,
		GRPCMethod:     msg.Method,
		URL:            msg.FullMethod,
		GRPCStreaming:  msg.IsStreaming,
		GRPCFrameIndex: msg.FrameIndex,
		GRPCCompressed: msg.Compressed,
	}

	if msg.JSON != "" {
		rec.GRPCData = msg.JSON
	} else if msg.Frame != nil {
		rec.Size = len(msg.Frame.Data)
	}

	if msg.Error != "" {
		rec.Error = msg.Error
		if msg.Frame != nil && len(msg.Frame.Data) > 0 {
			rec.GRPCRawData = base64.StdEncoding.EncodeToString(msg.Frame.Data)
		}
	}

	s.recorder.write(rec)
}

func isPrintableText(data []byte) bool {
	for _, b := range data {
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			return false
		}
		if b == 127 {
			return false
		}
	}
	return true
}

func (s *Session) Debug(format string, args ...interface{}) {
	if s.recorder.logLevel < LogLevelDebug {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "debug",
		Host:        s.Host,
		Error:       fmt.Sprintf(format, args...),
	}

	s.recorder.write(rec)
}

func (s *Session) LogError(err error) {
	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "error",
		Host:        s.Host,
		Error:       err.Error(),
	}

	s.recorder.write(rec)
}

func cloneHeaders(h http.Header) map[string][]string {
	if h == nil {
		return nil
	}
	clone := make(map[string][]string, len(h))
	for k, v := range h {
		clone[k] = append([]string(nil), v...)
	}
	return clone
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// SessionLogger wraps Session to implement Logger interface.
type SessionLogger struct {
	*Session
}

var _ Logger = (*SessionLogger)(nil)

func (s *Session) Logger() *SessionLogger {
	return &SessionLogger{Session: s}
}
