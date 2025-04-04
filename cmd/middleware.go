// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex service HTTP server middle-wares.

package cmd

import (
	"net/http"

	"github.com/rs/zerolog"
)

// verboseResponseWriter is a verbose wrapper around [http.ResponseWriter].
type verboseResponseWriter struct {
	http.ResponseWriter

	// statusCode exposes the status code written by the request handler.
	statusCode int
}

// WriteHeader wraps [http.ResponseWriter]'s WriteHeader function.
func (w *verboseResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// loggingHTTPMiddleware logs incoming requests to HTTP handlers and their
// responses. Usage example:
//
//	withMiddleware := loggingHTTPMiddleware(&logger)
//	newHandler := withMiddleware(handler)
func loggingHTTPMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			verboseWriter := verboseResponseWriter{ResponseWriter: w}
			handler.ServeHTTP(&verboseWriter, r)

			logger.Debug().
				Str("remoteaddr", r.RemoteAddr).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", verboseWriter.statusCode).
				Msg("Request")
		})
	}
}
