package config_test

import (
	"testing"

	"github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	want := &config.Config{
		HTTPPort:   8080,
		SOCKS5Port: 1080,
		APIPort:    9090,
		CertDir:    "~/.cursor-tap",
		DataDir:    "~/.cursor-tap/data",
	}

	got := config.DefaultConfig()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("DefaultConfig() mismatch (-want +got):\n%s", diff)
	}

	require.Equal(t, httpstream.LogLevel(0), got.HTTPLogLevel)
	require.Empty(t, got.HTTPRecordFile)
	require.False(t, got.EnableHTTPParsing)
}
