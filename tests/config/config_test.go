package config_test

import (
	"testing"

	"github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/httpstream"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.HTTPPort != 8080 {
		t.Errorf("HTTPPort = %d, want 8080", cfg.HTTPPort)
	}
	if cfg.SOCKS5Port != 1080 {
		t.Errorf("SOCKS5Port = %d, want 1080", cfg.SOCKS5Port)
	}
	if cfg.APIPort != 8888 {
		t.Errorf("APIPort = %d, want 8888", cfg.APIPort)
	}
	if cfg.CertDir != "~/.cursor-tap" {
		t.Errorf("CertDir = %q, want ~/.cursor-tap", cfg.CertDir)
	}
	if cfg.DataDir != "~/.cursor-tap/data" {
		t.Errorf("DataDir = %q, want ~/.cursor-tap/data", cfg.DataDir)
	}
	if cfg.EnableHTTPParsing {
		t.Error("EnableHTTPParsing should default to false")
	}
	if cfg.HTTPLogLevel != httpstream.LogLevel(0) {
		t.Errorf("HTTPLogLevel = %d, want 0", cfg.HTTPLogLevel)
	}
	if cfg.HTTPRecordFile != "" {
		t.Errorf("HTTPRecordFile = %q, want empty", cfg.HTTPRecordFile)
	}
}
