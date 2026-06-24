package httpstream_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
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
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ParseContentEncoding(%q) mismatch (-want +got):\n%s", tt.value, diff)
			}
		})
	}
}

func TestBodyDecoderGzip(t *testing.T) {
	plain := []byte("decoded body")
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	_, err := w.Write(plain)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	decoder := httpstream.NewBodyDecoder()
	headers := http.Header{"Content-Encoding": []string{"gzip"}}
	reader := decoder.Decode(bytes.NewReader(compressed.Bytes()), headers)

	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, plain, got)
}

func TestBodyDecoderIdentitySkipped(t *testing.T) {
	plain := []byte("plain")
	decoder := httpstream.NewBodyDecoder()
	headers := http.Header{"Content-Encoding": []string{"identity"}}
	reader := decoder.Decode(bytes.NewReader(plain), headers)

	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, plain, got)
}
