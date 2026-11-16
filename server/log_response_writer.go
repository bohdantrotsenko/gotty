package server

import (
	"bufio"
	"net"
	"net/http"
)

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, _ := w.ResponseWriter.(http.Hijacker)
	w.status = http.StatusSwitchingProtocols
	return hj.Hijack()
}

// Unwrap returns the underlying ResponseWriter so that http.ResponseController
// can reach the wrapped writer's capabilities (e.g. Flush, Hijack).
func (w *logResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
