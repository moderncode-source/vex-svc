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

// HealthHandler handles /healthz server endpoint. Used as a liveness probe
// during deployment to check whether the server is even running.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
