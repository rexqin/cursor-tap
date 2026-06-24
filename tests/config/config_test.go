package config_test

import (
	"path/filepath"
	"testing"

	"github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestDefaultProxyConfig(t *testing.T) {
	certDir := "~/.cursor-tap"
	want := &config.ProxyConfig{
		HTTPPort:     8080,
		SOCKS5Port:   1080,
		CertDir:      certDir,
		DataDir:      filepath.Join(certDir, "data"),
		RecordDB:     filepath.Join(certDir, "data", "records.db"),
		APINotifyURL: "http://127.0.0.1:9090/internal/notify",
	}

	got := config.DefaultProxyConfig()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("DefaultProxyConfig() mismatch (-want +got):\n%s", diff)
	}

	require.Equal(t, httpstream.LogLevel(0), got.HTTPLogLevel)
	require.False(t, got.EnableHTTPParsing)
}

func TestDefaultAPIConfig(t *testing.T) {
	certDir := "~/.cursor-tap"
	want := &config.APIConfig{
		Port:     9090,
		CertDir:  certDir,
		RecordDB: filepath.Join(certDir, "data", "records.db"),
	}

	got := config.DefaultAPIConfig()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("DefaultAPIConfig() mismatch (-want +got):\n%s", diff)
	}
}
