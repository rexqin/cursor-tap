package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	appconfig "github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/apiserver"
	"github.com/spf13/cobra"
)

var (
	apiPort  int
	certDir  string
	recordDB string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tap-api",
		Short: "Cursor-Tap management API server",
		Long:  `Standalone REST and WebSocket API for the Cursor-Tap Inspector.`,
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the API server",
		RunE:  runStart,
	}
	startCmd.Flags().IntVar(&apiPort, "port", 9090, "API server port")
	startCmd.Flags().StringVar(&certDir, "cert-dir", "~/.cursor-tap", "Certificate storage directory")
	startCmd.Flags().StringVar(&recordDB, "record-db", "", "SQLite record database (default: cert-dir/data/records.db)")

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runStart(cmd *cobra.Command, args []string) error {
	certDir = expandPath(certDir)
	if recordDB == "" {
		recordDB = filepath.Join(certDir, "data", "records.db")
	} else {
		recordDB = expandPath(recordDB)
	}

	config := appconfig.APIConfig{
		Port:     apiPort,
		CertDir:  certDir,
		RecordDB: recordDB,
	}

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         cursor-tap API Starting          ║")
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Printf("║  API Server:    127.0.0.1:%-15d║\n", config.Port)
	fmt.Printf("║  Cert Dir:      %-25s║\n", truncateString(config.CertDir, 25))
	fmt.Printf("║  Record DB:     %-25s║\n", truncateString(config.RecordDB, 25))
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop...")
	fmt.Println()

	server, err := apiserver.New(config)
	if err != nil {
		return fmt.Errorf("create API server: %w", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		server.Stop()
	}()

	return server.Start()
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
