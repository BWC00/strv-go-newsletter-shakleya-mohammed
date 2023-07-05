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


// Package requestlog provides utility functions and types for logging HTTP request information.
package requestlog

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

// logEntry represents a log entry for an HTTP request.
type logEntry struct {
	ReceivedTime      time.Time	// The time the request was received
	RequestMethod     string	// The HTTP method of the request (e.g., GET, POST)
	RequestURL        string	// The URL of the request
	RequestHeaderSize int64		// The size of the request headers
	RequestBodySize   int64		// The size of the request body
	UserAgent         string	// The user agent of the client
	Referer           string	// The referer URL of the client
	Proto             string	// The protocol version of the request

	RemoteIP string	// The IP address of the client
	ServerIP string	// The IP address of the server

	Status             int			 // The HTTP response status code
	ResponseHeaderSize int64		 // The size of the response headers
	ResponseBodySize   int64		 // The size of the response body
	Latency            time.Duration // The duration of the request processing
}

// ipFromHostPort extracts the IP address from a host:port string.
func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}

// readCounterCloser is a wrapper around an io.ReadCloser that keeps track of the number of bytes read.
type readCounterCloser struct {
	r   io.ReadCloser
	n   int64
	err error
}

// Read reads bytes from the underlying reader and updates the byte count.
func (rcc *readCounterCloser) Read(p []byte) (n int, err error) {
	if rcc.err != nil {
		return 0, rcc.err
	}
	n, rcc.err = rcc.r.Read(p)
	rcc.n += int64(n)
	return n, rcc.err
}

// Close closes the underlying reader and returns an error indicating that the reader is closed.
func (rcc *readCounterCloser) Close() error {
	rcc.err = errors.New("read from closed reader")
	return rcc.r.Close()
}

// writeCounter is an integer type that implements the io.Writer interface.
type writeCounter int64

// Write increments the write counter by the number of bytes written.
func (wc *writeCounter) Write(p []byte) (n int, err error) {
	*wc += writeCounter(len(p))
	return len(p), nil
}

// headerSize calculates the size of the headers in bytes.
func headerSize(h http.Header) int64 {
	var wc writeCounter
	h.Write(&wc)
	return int64(wc) + 2 // Add 2 bytes for CRLF (carriage return, line feed)
}

// responseStats is a wrapper around an http.ResponseWriter that captures response information.
type responseStats struct {
	w     http.ResponseWriter // The underlying ResponseWriter
	hsize int64				  // The size of the response headers
	wc    writeCounter		  // The write counter for the response body
	code  int				  // The HTTP response status code
}

// Header returns the header map of the underlying ResponseWriter.
func (r *responseStats) Header() http.Header {
	return r.w.Header()
}

// WriteHeader writes the HTTP response status code and captures the header size.
func (r *responseStats) WriteHeader(statusCode int) {
	if r.code != 0 {
		return
	}
	r.hsize = headerSize(r.w.Header())
	r.w.WriteHeader(statusCode)
	r.code = statusCode
}

// Write writes the response body and updates the write counter.
func (r *responseStats) Write(p []byte) (n int, err error) {
	if r.code == 0 {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.w.Write(p)
	r.wc.Write(p[:n])
	return
}

// size returns the header size and response body size.
func (r *responseStats) size() (hdr, body int64) {
	if r.code == 0 {
		return headerSize(r.w.Header()), 0
	}
	// Use the header size from the time WriteHeader was called.
	// The Header map can be mutated after the call to add HTTP Trailers,
	// which we don't want to count.
	return r.hsize, int64(r.wc)
}
