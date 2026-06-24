package httpstream_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
)

func TestParseMethodFromURL(t *testing.T) {
	service, method, full := httpstream.ParseMethodFromURL("/aiserver.v1.AiService/RunSSE")
	if service != "aiserver.v1.AiService" {
		t.Errorf("service = %q", service)
	}
	if method != "RunSSE" {
		t.Errorf("method = %q", method)
	}
	if full != "/aiserver.v1.AiService/RunSSE" {
		t.Errorf("fullMethod = %q", full)
	}

	service, method, full = httpstream.ParseMethodFromURL("invalid")
	if service != "" || method != "" {
		t.Errorf("invalid path should return empty service/method, got %q/%q", service, method)
	}
	if full != "invalid" {
		t.Errorf("fullMethod = %q, want invalid", full)
	}
}

func TestGRPCParserReadFrame(t *testing.T) {
	parser := httpstream.NewGRPCParser(nil)

	payload := []byte("hello")
	var buf bytes.Buffer
	buf.WriteByte(0)
	buf.Write([]byte{0, 0, 0, byte(len(payload))})
	buf.Write(payload)

	frame, err := parser.ReadFrame(&buf)
	if err != nil {
		t.Fatalf("ReadFrame: %v", err)
	}
	if frame.Compressed {
		t.Fatal("expected uncompressed frame")
	}
	if string(frame.Data) != "hello" {
		t.Fatalf("Data = %q, want hello", frame.Data)
	}
}

func TestGRPCParserReadFrameGzip(t *testing.T) {
	parser := httpstream.NewGRPCParser(nil)

	plain := []byte(`{"msg":"test"}`)
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	if _, err := w.Write(plain); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	compressed := gz.Bytes()

	var buf bytes.Buffer
	buf.WriteByte(1)
	length := make([]byte, 4)
	length[3] = byte(len(compressed))
	buf.Write(length)
	buf.Write(compressed)

	frame, err := parser.ReadFrame(&buf)
	if err != nil {
		t.Fatalf("ReadFrame: %v", err)
	}
	if !frame.Compressed {
		t.Fatal("expected compressed frame")
	}
	if string(frame.Data) != string(plain) {
		t.Fatalf("Data = %q, want %q", frame.Data, plain)
	}
}

func TestGRPCParserReadAllFramesEOF(t *testing.T) {
	parser := httpstream.NewGRPCParser(nil)
	frames, err := parser.ReadAllFrames(bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("ReadAllFrames: %v", err)
	}
	if len(frames) != 0 {
		t.Fatalf("expected 0 frames, got %d", len(frames))
	}

	payload := []byte("x")
	var buf bytes.Buffer
	buf.WriteByte(0)
	buf.Write([]byte{0, 0, 0, 1})
	buf.Write(payload)

	frames, err = parser.ReadAllFrames(&buf)
	if err != nil && err != io.EOF {
		t.Fatalf("ReadAllFrames: %v", err)
	}
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
}
