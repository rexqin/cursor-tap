package httpstream_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/stretchr/testify/require"
)

func TestParseMethodFromURL(t *testing.T) {
	t.Run("valid path", func(t *testing.T) {
		service, method, full := httpstream.ParseMethodFromURL("/aiserver.v1.AiService/RunSSE")
		require.Equal(t, "aiserver.v1.AiService", service)
		require.Equal(t, "RunSSE", method)
		require.Equal(t, "/aiserver.v1.AiService/RunSSE", full)
	})

	t.Run("invalid path", func(t *testing.T) {
		service, method, full := httpstream.ParseMethodFromURL("invalid")
		require.Empty(t, service)
		require.Empty(t, method)
		require.Equal(t, "invalid", full)
	})
}

func TestGRPCParserReadFrame(t *testing.T) {
	parser := httpstream.NewGRPCParser(nil)

	payload := []byte("hello")
	var buf bytes.Buffer
	buf.WriteByte(0)
	buf.Write([]byte{0, 0, 0, byte(len(payload))})
	buf.Write(payload)

	frame, err := parser.ReadFrame(&buf)
	require.NoError(t, err)
	require.False(t, frame.Compressed)
	require.Equal(t, []byte("hello"), frame.Data)
}

func TestGRPCParserReadFrameGzip(t *testing.T) {
	parser := httpstream.NewGRPCParser(nil)

	plain := []byte(`{"msg":"test"}`)
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	_, err := w.Write(plain)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	compressed := gz.Bytes()

	var buf bytes.Buffer
	buf.WriteByte(1)
	length := make([]byte, 4)
	length[3] = byte(len(compressed))
	buf.Write(length)
	buf.Write(compressed)

	frame, err := parser.ReadFrame(&buf)
	require.NoError(t, err)
	require.True(t, frame.Compressed)
	require.Equal(t, plain, frame.Data)
}

func TestGRPCParserReadAllFramesEOF(t *testing.T) {
	parser := httpstream.NewGRPCParser(nil)

	frames, err := parser.ReadAllFrames(bytes.NewReader(nil))
	require.NoError(t, err)
	require.Empty(t, frames)

	payload := []byte("x")
	var buf bytes.Buffer
	buf.WriteByte(0)
	buf.Write([]byte{0, 0, 0, 1})
	buf.Write(payload)

	frames, err = parser.ReadAllFrames(&buf)
	if err != nil && err != io.EOF {
		require.NoError(t, err)
	}
	require.Len(t, frames, 1)
}
