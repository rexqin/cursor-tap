package proxy_test

import (
	"path/filepath"
	"testing"

	appconfig "github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/proxy"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	dir := t.TempDir()
	cfg := appconfig.ProxyConfig{
		HTTPPort:          18080,
		SOCKS5Port:        11080,
		CertDir:           dir,
		DataDir:           filepath.Join(dir, "data"),
		RecordDB:          filepath.Join(dir, "data", "records.db"),
		APINotifyURL:      "http://127.0.0.1:9090/internal/notify",
		EnableHTTPParsing: true,
	}

	server, err := proxy.NewServer(cfg)
	require.NoError(t, err)
	t.Cleanup(server.Close)
}
