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

func loggingHTTPMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug().
				Str("remoteaddr", r.RemoteAddr).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("request")
			handler.ServeHTTP(w, r)
		})
	}
}
