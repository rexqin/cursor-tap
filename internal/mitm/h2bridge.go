package mitm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/burpheart/cursor-tap/internal/httpstream"
	"golang.org/x/net/http2"
)

// alpnProtos lists protocols offered during TLS MITM handshakes.
// Cursor clients may require HTTP/2 (h2) exclusively.
func alpnProtos() []string {
	return []string{"h2", "http/1.1"}
}

// isH2 reports whether the negotiated ALPN protocol is HTTP/2.
func isH2(proto string) bool {
	return proto == "h2"
}

// pipeClientH2 handles a client connection that negotiated HTTP/2.
func (i *Interceptor) pipeClientH2(client, server net.Conn, host, serverProto string) error {
	if !i.enableHTTPParsing && isH2(serverProto) {
		return i.pipeSimple(client, server)
	}

	bridge := &h2Bridge{
		interceptor: i,
		server:      server,
		serverProto: serverProto,
		host:        host,
	}

	if i.enableHTTPParsing {
		bridge.parser = i.newHTTPParser(host)
	}

	h2srv := &http2.Server{}
	fmt.Printf("[DEBUG] Serving HTTP/2 for %s (server ALPN: %q)\n", host, serverProto)
	h2srv.ServeConn(client, &http2.ServeConnOpts{
		Context: context.Background(),
		Handler: http.HandlerFunc(bridge.handle),
	})
	return nil
}

type h2Bridge struct {
	interceptor *Interceptor
	server      net.Conn
	serverProto string
	host        string
	parser      *httpstream.Parser
	h1Mu    sync.Mutex
	h2cc    *http2.ClientConn
	h2ccErr error
	h2once  sync.Once
}

func (b *h2Bridge) handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host == "" {
		r.URL.Host = b.host
	}
	if r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}
	r.Host = b.host

	if b.parser != nil {
		b.parser.ProcessRequest(r)
	}

	resp, err := b.roundTrip(r)
	if err != nil {
		fmt.Printf("[DEBUG] H2 proxy round-trip error for %s %s: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if b.parser != nil {
		b.parser.ProcessResponse(resp)
	}

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (b *h2Bridge) roundTrip(r *http.Request) (*http.Response, error) {
	if isH2(b.serverProto) {
		cc, err := b.getH2Client()
		if err != nil {
			return nil, err
		}
		return cc.RoundTrip(r)
	}
	return b.roundTripH1(r)
}

func (b *h2Bridge) getH2Client() (*http2.ClientConn, error) {
	b.h2once.Do(func() {
		tr := &http2.Transport{}
		b.h2cc, b.h2ccErr = tr.NewClientConn(b.server)
	})
	return b.h2cc, b.h2ccErr
}

func (b *h2Bridge) roundTripH1(r *http.Request) (*http.Response, error) {
	b.h1Mu.Lock()
	defer b.h1Mu.Unlock()

	req := r.Clone(r.Context())
	req.RequestURI = ""

	if err := req.Write(b.server); err != nil {
		return nil, err
	}
	return http.ReadResponse(bufio.NewReader(b.server), req)
}
