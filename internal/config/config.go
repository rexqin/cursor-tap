// Package config holds application configuration for the MITM proxy.
package config

import "github.com/burpheart/cursor-tap/internal/httpstream"

// Config holds the application configuration.
type Config struct {
	HTTPPort      int    `json:"http_port"`
	SOCKS5Port    int    `json:"socks5_port"`
	APIPort       int    `json:"api_port"`
	CertDir       string `json:"cert_dir"`
	DataDir       string `json:"data_dir"`
	UpstreamProxy string `json:"upstream_proxy"` // e.g., "http://127.0.0.1:7890" or "socks5://127.0.0.1:1080"

	// HTTP parsing options
	EnableHTTPParsing bool                `json:"enable_http_parsing"`
	HTTPLogLevel      httpstream.LogLevel `json:"http_log_level"`
	HTTPRecordFile    string              `json:"http_record_file"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		HTTPPort:   8080,
		SOCKS5Port: 1080,
		APIPort:    9090,
		CertDir:    "~/.cursor-tap",
		DataDir:    "~/.cursor-tap/data",
	}
}
