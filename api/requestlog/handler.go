// Copyright 2019 The Go Cloud Development Kit Authors
// Modified by Mohammed Shakleya
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Source: https://github.com/google/go-cloud/blob/master/server/requestlog/requestlog.go


// Package requestlog provides an HTTP handler for logging request information.
// It logs details such as the request method, URL, headers, user agent, response status, latency, and more.
// The logged information is sent to a logger instance.
package requestlog

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
)

// Handler is an HTTP handler that logs request information.
type Handler struct {
	handler http.Handler
	logger  *logger.Logger
}

// NewHandler creates a new request log handler with the given HTTP handler and logger.
func NewHandler(h http.HandlerFunc, l *logger.Logger) *Handler {
	return &Handler{
		handler: h,
		logger:  l,
	}
}

// ServeHTTP implements the http.Handler interface.
// It logs request information such as method, URL, headers, user agent, response status, latency, and more.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Create a log entry to store request details
	le := &logEntry{
		ReceivedTime:      start,
		RequestMethod:     r.Method,
		RequestURL:        r.URL.String(),
		RequestHeaderSize: headerSize(r.Header),
		UserAgent:         r.UserAgent(),
		Referer:           r.Referer(),
		Proto:             r.Proto,
		RemoteIP:          ipFromHostPort(r.RemoteAddr),
	}

	// Retrieve the server IP address from the request context if available
	if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		le.ServerIP = ipFromHostPort(addr.String())
	}

	// Create a new request and response writer to capture response information
	r2 := new(http.Request)
	*r2 = *r
	rcc := &readCounterCloser{r: r.Body}
	r2.Body = rcc
	w2 := &responseStats{w: w}

	// Call the underlying handler with the modified request and response writer
	h.handler.ServeHTTP(w2, r2)

	// Update the log entry with response details
	le.Latency = time.Since(start)
	if rcc.err == nil && rcc.r != nil {
		// If the handler hasn't encountered an error in the Body (like EOF),
		// then consume the rest of the Body to provide an accurate rcc.n.
		io.Copy(io.Discard, rcc)
	}
	le.RequestBodySize = rcc.n
	le.Status = w2.code
	if le.Status == 0 {
		le.Status = http.StatusOK
	}
	le.ResponseHeaderSize, le.ResponseBodySize = w2.size()

	// Log the request information using the provided logger
	h.logger.Info().
		Time("received_time", le.ReceivedTime).
		Str("method", le.RequestMethod).
		Str("url", le.RequestURL).
		Int64("header_size", le.RequestHeaderSize).
		Int64("body_size", le.RequestBodySize).
		Str("agent", le.UserAgent).
		Str("referer", le.Referer).
		Str("proto", le.Proto).
		Str("remote_ip", le.RemoteIP).
		Str("server_ip", le.ServerIP).
		Int("status", le.Status).
		Int64("resp_header_size", le.ResponseHeaderSize).
		Int64("resp_body_size", le.ResponseBodySize).
		Dur("latency", le.Latency).
		Msg("")
}
