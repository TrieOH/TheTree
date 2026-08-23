package store

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// sseKeepaliveIntervalNS is how often the stream emits a `: ping` comment
// line. Atomic so tests can shorten it without racing in-flight streams
// (a previous test's stream goroutine may still be reading it).
var sseKeepaliveIntervalNS atomic.Int64

func init() { sseKeepaliveIntervalNS.Store(int64(30 * time.Second)) }

// sseStream is a hijacked close-delimited HTTP/1.1 response speaking SSE.
// The harness's middleware stack wraps the ResponseWriter in writers that
// expose no Flush and no write-deadline control, and its http.Server sets a
// 60s WriteTimeout — all fatal for a long-lived stream. Hijacking the conn
// clears every deadline (net/http does that on Hijack) and hands the raw
// connection to us: headers are written once, then each event is a framed
// write + flush. No Content-Length / Transfer-Encoding means the body is
// close-delimited — legal HTTP/1.1 and exactly what SSE wants.
//
// Only one goroutine may write at a time; event() takes the mutex.
type sseStream struct {
	conn net.Conn
	bw   *bufio.Writer
	mu   sync.Mutex

	done      chan struct{}
	closeOnce sync.Once
}

// corsResponseHeaders are the response headers the harness CORS middleware
// (rs/cors with AllowCredentials) sets on the ResponseWriter for an actual
// request. Hijack discards the ResponseWriter, so the hand-built head must
// re-emit them or the browser blocks the cross-origin stream (EventSource
// is CORS-gated). Connection/Content-Length/Transfer-Encoding are hop-by-hop
// and deliberately excluded: the body is close-delimited.
var corsResponseHeaders = []string{
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Credentials",
	"Access-Control-Expose-Headers",
	"Vary",
}

// newSSEStream hijacks the connection and writes the response head. Any
// error means the connection was not taken over (the caller can still write
// a plain HTTP error).
func newSSEStream(w http.ResponseWriter) (*sseStream, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("response writer does not support hijacking")
	}
	conn, rw, err := hj.Hijack()
	if err != nil {
		return nil, err
	}
	bw := rw.Writer

	// Snapshot the CORS headers the middleware set before hijacking (the
	// ResponseWriter is discarded on Hijack) and splice them into the raw
	// head; no Origin header means none were set and none are emitted.
	head := strings.Builder{}
	head.WriteString("HTTP/1.1 200 OK\r\n")
	for _, name := range corsResponseHeaders {
		for _, v := range w.Header().Values(name) {
			head.WriteString(name + ": " + v + "\r\n")
		}
	}
	head.WriteString("Content-Type: text/event-stream; charset=utf-8\r\n" +
		"Cache-Control: no-cache\r\n" +
		"Connection: keep-alive\r\n" +
		"X-Accel-Buffering: no\r\n" +
		"\r\n")
	_, err = bw.WriteString(head.String())
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	err = bw.Flush()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &sseStream{conn: conn, bw: bw, done: make(chan struct{})}, nil
}

// event writes one SSE event block (`event: <name>\ndata: <json>\n\n`) and
// flushes it to the wire. Returns the write error (the caller treats a dead
// conn as terminal).
func (s *sseStream) event(name string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := fmt.Fprintf(s.bw, "event: %s\n", name)
	if err != nil {
		return err
	}
	// SSE forbids bare newlines inside data; JSON here is single-line, but
	// write defensively so a multi-line payload cannot desync the protocol.
	for _, line := range splitLines(data) {
		_, err := fmt.Fprintf(s.bw, "data: %s\n", line)
		if err != nil {
			return err
		}
	}
	_, err = s.bw.WriteString("\n")
	if err != nil {
		return err
	}
	return s.bw.Flush()
}

// keepalive writes a comment line (`: ping`) every ~30s so proxies and the
// client can tell the stream is alive. Exits when the stream closes.
func (s *sseStream) keepalive() {
	t := time.NewTicker(time.Duration(sseKeepaliveIntervalNS.Load()))
	defer t.Stop()
	for {
		select {
		case <-t.C:
			s.mu.Lock()
			_, _ = s.bw.WriteString(": ping\n\n")
			err := s.bw.Flush()
			s.mu.Unlock()
			if err != nil {
				s.close() // conn died under us — tear down
				return
			}
		case <-s.done:
			return
		}
	}
}

// waitClosed blocks until the client goes away (the client never sends
// data, so the first read fails only on disconnect), then closes the
// stream — the close-delimited body terminator.
func (s *sseStream) waitClosed() {
	buf := make([]byte, 1)
	_, _ = s.conn.Read(buf)
	s.close()
}

// close tears the connection down (idempotent).
func (s *sseStream) close() {
	s.closeOnce.Do(func() {
		close(s.done)
		_ = s.conn.Close()
	})
}

// splitLines splits data into lines (SSE forbids bare newlines in data).
func splitLines(b []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, c := range b {
		if c == '\n' {
			lines = append(lines, b[start:i])
			start = i + 1
		}
	}
	return append(lines, b[start:])
}
