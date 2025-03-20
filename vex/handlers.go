// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex HTTP request handlers.

package vex

import "net/http"

// HealthHandler handles requests to service liveness probe endpoint that can
// be used to check whether the server is running.
func HealthHandler(w http.ResponseWriter, req *http.Request) {
	if len(req.Method) != 0 && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ReadyHandler handles requests to service readiness probe endpoint that can
// be used to check whether the server is ready to receive traffic.
func ReadyHandler(w http.ResponseWriter, req *http.Request) {
	if len(req.Method) != 0 && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}
