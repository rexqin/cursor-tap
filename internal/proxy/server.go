// Package proxy provides HTTP and SOCKS5 proxy servers with TLS MITM support.
package proxy

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/burpheart/cursor-tap/internal/ca"
	appconfig "github.com/burpheart/cursor-tap/internal/config"
	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/burpheart/cursor-tap/internal/mitm"
	"github.com/burpheart/cursor-tap/internal/notify"
)

// Server is the MITM proxy server (HTTP + SOCKS5).
type Server struct {
	config      appconfig.ProxyConfig
	ca          *ca.CA
	interceptor *mitm.Interceptor
	keyLog      *mitm.KeyLogWriter
	recorder    *httpstream.Recorder
	notifier    *notify.Client

	httpListener   net.Listener
	socks5Listener net.Listener

	mu       sync.Mutex
	running  bool
	stopChan chan struct{}
}

// NewServer creates a new proxy server.
func NewServer(config appconfig.ProxyConfig) (*Server, error) {
	if err := os.MkdirAll(config.CertDir, 0755); err != nil {
		return nil, fmt.Errorf("create cert dir: %w", err)
	}
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	caInstance, err := ca.New(ca.Options{CertDir: config.CertDir})
	if err != nil {
		return nil, fmt.Errorf("initialize CA: %w", err)
	}

	keyLogPath := filepath.Join(config.DataDir, "sslkeys.log")
	keyLog, err := mitm.NewKeyLogWriter(keyLogPath)
	if err != nil {
		return nil, fmt.Errorf("create keylog writer: %w", err)
	}
	fmt.Printf("[INFO] TLS KeyLog enabled: %s\n", keyLogPath)

	var interceptorOpts []mitm.InterceptorOption
	var recorder *httpstream.Recorder
	notifier := notify.NewClient(config.APINotifyURL)

	if config.EnableHTTPParsing {
		interceptorOpts = append(interceptorOpts, mitm.WithHTTPParsing(true))

		httpLogLevel := config.HTTPLogLevel
		httpLogger := httpstream.NewDefaultLogger(
			httpstream.WithLevel(httpLogLevel),
			httpstream.WithColor(true),
		)
		interceptorOpts = append(interceptorOpts, mitm.WithHTTPLogger(httpLogger))
		fmt.Printf("[INFO] HTTP parsing enabled (log level: %d)\n", config.HTTPLogLevel)

		recorder, err = httpstream.NewRecorder(
			config.RecordDB,
			httpstream.WithRecorderLogLevel(httpLogLevel),
			httpstream.WithOnRecord(func(id int64, _ httpstream.Record) {
				notifier.NotifyLatest(id)
			}),
		)
		if err != nil {
			keyLog.Close()
			return nil, fmt.Errorf("create HTTP recorder: %w", err)
		}
		interceptorOpts = append(interceptorOpts, mitm.WithRecorder(recorder))
		fmt.Printf("[INFO] HTTP recording enabled: %s\n", config.RecordDB)
		fmt.Printf("[INFO] API notify URL: %s\n", config.APINotifyURL)
	}

	interceptor := mitm.NewInterceptor(caInstance, keyLog, config.UpstreamProxy, interceptorOpts...)

	return &Server{
		config:      config,
		ca:          caInstance,
		interceptor: interceptor,
		keyLog:      keyLog,
		recorder:    recorder,
		notifier:    notifier,
		stopChan:    make(chan struct{}),
	}, nil
}

// Start starts HTTP and SOCKS5 proxy servers.
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("server already running")
	}
	s.running = true
	s.mu.Unlock()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.startHTTPProxy(); err != nil {
			errChan <- fmt.Errorf("HTTP proxy: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.startSOCKS5Proxy(); err != nil {
			errChan <- fmt.Errorf("SOCKS5 proxy: %w", err)
		}
	}()

	select {
	case <-s.stopChan:
	case err := <-errChan:
		return err
	}

	wg.Wait()
	return nil
}

// Stop stops proxy servers.
func (s *Server) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopChan)

	if s.httpListener != nil {
		s.httpListener.Close()
	}
	if s.socks5Listener != nil {
		s.socks5Listener.Close()
	}
	if s.keyLog != nil {
		s.keyLog.Close()
	}
	if s.recorder != nil {
		s.recorder.Close()
	}
}

// Close releases file handles. Safe without Start.
func (s *Server) Close() {
	if s.keyLog != nil {
		s.keyLog.Close()
		s.keyLog = nil
	}
	if s.recorder != nil {
		s.recorder.Close()
		s.recorder = nil
	}
}

// startHTTPProxy starts the HTTP/HTTPS proxy server.
func (s *Server) startHTTPProxy() error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.config.HTTPPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	s.httpListener = listener
	fmt.Printf("[INFO] HTTP proxy listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.stopChan:
				return nil
			default:
				fmt.Printf("[ERROR] HTTP accept: %v\n", err)
				continue
			}
		}

		go s.handleHTTPConnection(conn)
	}
}

func (s *Server) handleHTTPConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	req, err := http.ReadRequest(reader)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("[DEBUG] HTTP read request error: %v\n", err)
		}
		return
	}

	if req.Method == http.MethodConnect {
		buffConn := &bufferedNetConn{Conn: conn, reader: reader}
		s.handleHTTPConnect(buffConn, req)
	} else {
		s.handleHTTPRequest(conn, req, reader)
	}
}

func (s *Server) handleHTTPConnect(conn net.Conn, req *http.Request) {
	host := req.Host
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	targetHost, targetPortStr, err := net.SplitHostPort(host)
	if err != nil {
		fmt.Printf("[ERROR] Invalid CONNECT host: %s\n", host)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	targetPort, _ := strconv.Atoi(targetPortStr)

	_, err = conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		return
	}

	fmt.Printf("[INFO] CONNECT %s:%d\n", targetHost, targetPort)

	if err := s.interceptor.InterceptAuto(conn, targetHost, targetPort); err != nil {
		if !isConnectionClosed(err) {
			fmt.Printf("[ERROR] Intercept %s:%d: %v\n", targetHost, targetPort, err)
		}
	}
}

func (s *Server) handleHTTPRequest(clientConn net.Conn, req *http.Request, _ *bufio.Reader) {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	if host == "" {
		clientConn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	if !strings.Contains(host, ":") {
		host = host + ":80"
	}

	targetConn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer targetConn.Close()

	req.URL.Scheme = ""
	req.URL.Host = ""
	req.RequestURI = req.URL.RequestURI()

	if err := req.Write(targetConn); err != nil {
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(targetConn), req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	resp.Write(clientConn)
}

func (s *Server) startSOCKS5Proxy() error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.config.SOCKS5Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	s.socks5Listener = listener
	fmt.Printf("[INFO] SOCKS5 proxy listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.stopChan:
				return nil
			default:
				fmt.Printf("[ERROR] SOCKS5 accept: %v\n", err)
				continue
			}
		}

		go s.handleSOCKS5Connection(conn)
	}
}

func (s *Server) handleSOCKS5Connection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	header := make([]byte, 2)
	if _, err := io.ReadFull(reader, header); err != nil {
		return
	}

	if header[0] != 0x05 {
		return
	}

	nmethods := int(header[1])
	methods := make([]byte, nmethods)
	if _, err := io.ReadFull(reader, methods); err != nil {
		return
	}

	conn.Write([]byte{0x05, 0x00})

	reqHeader := make([]byte, 4)
	if _, err := io.ReadFull(reader, reqHeader); err != nil {
		return
	}

	if reqHeader[0] != 0x05 || reqHeader[1] != 0x01 {
		conn.Write([]byte{0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}

	var targetHost string
	var targetPort int

	switch reqHeader[3] {
	case 0x01:
		addr := make([]byte, 4)
		if _, err := io.ReadFull(reader, addr); err != nil {
			return
		}
		targetHost = net.IP(addr).String()
	case 0x03:
		lenByte := make([]byte, 1)
		if _, err := io.ReadFull(reader, lenByte); err != nil {
			return
		}
		domain := make([]byte, lenByte[0])
		if _, err := io.ReadFull(reader, domain); err != nil {
			return
		}
		targetHost = string(domain)
	case 0x04:
		addr := make([]byte, 16)
		if _, err := io.ReadFull(reader, addr); err != nil {
			return
		}
		targetHost = net.IP(addr).String()
	default:
		conn.Write([]byte{0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}

	portBytes := make([]byte, 2)
	if _, err := io.ReadFull(reader, portBytes); err != nil {
		return
	}
	targetPort = int(binary.BigEndian.Uint16(portBytes))

	fmt.Printf("[INFO] SOCKS5 CONNECT %s:%d\n", targetHost, targetPort)

	response := []byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
	conn.Write(response)

	conn.SetDeadline(time.Time{})

	buffConn := &bufferedNetConn{Conn: conn, reader: reader}

	if err := s.interceptor.InterceptAuto(buffConn, targetHost, targetPort); err != nil {
		if !isConnectionClosed(err) {
			fmt.Printf("[ERROR] SOCKS5 intercept %s:%d: %v\n", targetHost, targetPort, err)
		}
	}
}

func isConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "EOF") ||
		errors.Is(err, io.EOF)
}

type bufferedNetConn struct {
	net.Conn
	reader *bufio.Reader
}

func (c *bufferedNetConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}
