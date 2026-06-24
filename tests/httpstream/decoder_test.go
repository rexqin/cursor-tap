package httpstream_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
)

func TestParseContentEncoding(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{"empty", "", nil},
		{"gzip", "gzip", []string{"gzip"}},
		{"multiple", "gzip, br", []string{"gzip", "br"}},
		{"identity skipped", "identity, gzip", []string{"gzip"}},
		{"whitespace", "  deflate , gzip  ", []string{"deflate", "gzip"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := httpstream.ParseContentEncoding(tt.value)
			if len(got) != len(tt.want) {
				t.Fatalf("ParseContentEncoding(%q) = %v, want %v", tt.value, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("ParseContentEncoding(%q) = %v, want %v", tt.value, got, tt.want)
				}
			}
		})
	}
}

func TestBodyDecoderGzip(t *testing.T) {
	plain := []byte("decoded body")
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	if _, err := w.Write(plain); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	decoder := httpstream.NewBodyDecoder()
	headers := http.Header{"Content-Encoding": []string{"gzip"}}
	reader := decoder.Decode(bytes.NewReader(compressed.Bytes()), headers)

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != string(plain) {
		t.Fatalf("decoded = %q, want %q", got, plain)
	}
}

func TestBodyDecoderIdentitySkipped(t *testing.T) {
	plain := []byte("plain")
	decoder := httpstream.NewBodyDecoder()
	headers := http.Header{"Content-Encoding": []string{"identity"}}
	reader := decoder.Decode(bytes.NewReader(plain), headers)

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != string(plain) {
		t.Fatalf("decoded = %q, want %q", got, plain)
	}
}
