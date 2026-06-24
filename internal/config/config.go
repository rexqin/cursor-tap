// Package config holds application configuration for the MITM proxy and API server.
package config

import (
	"path/filepath"

	"github.com/burpheart/cursor-tap/internal/httpstream"
)

// ProxyConfig holds tap proxy configuration.
type ProxyConfig struct {
	HTTPPort      int    `json:"http_port"`
	SOCKS5Port    int    `json:"socks5_port"`
	CertDir       string `json:"cert_dir"`
	DataDir       string `json:"data_dir"`
	UpstreamProxy string `json:"upstream_proxy"`
	RecordDB      string `json:"record_db"`
	APINotifyURL  string `json:"api_notify_url"`

	EnableHTTPParsing bool                `json:"enable_http_parsing"`
	HTTPLogLevel      httpstream.LogLevel `json:"http_log_level"`
}

// DefaultProxyConfig returns default proxy configuration.
func DefaultProxyConfig() *ProxyConfig {
	certDir := "~/.cursor-tap"
	dataDir := filepath.Join(certDir, "data")
	return &ProxyConfig{
		HTTPPort:     8080,
		SOCKS5Port:   1080,
		CertDir:      certDir,
		DataDir:      dataDir,
		RecordDB:     filepath.Join(dataDir, "records.db"),
		APINotifyURL: "http://127.0.0.1:9090/internal/notify",
	}
}

// APIConfig holds tap-api server configuration.
type APIConfig struct {
	Port     int    `json:"port"`
	CertDir  string `json:"cert_dir"`
	RecordDB string `json:"record_db"`
}

// DefaultAPIConfig returns default API server configuration.
func DefaultAPIConfig() *APIConfig {
	certDir := "~/.cursor-tap"
	dataDir := filepath.Join(certDir, "data")
	return &APIConfig{
		Port:     9090,
		CertDir:  certDir,
		RecordDB: filepath.Join(dataDir, "records.db"),
	}
}
