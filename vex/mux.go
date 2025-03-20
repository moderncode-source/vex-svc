// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex HTTP request multiplexer. See [http.ServeMux].

package vex

import "net/http"

// ServiceMux is the default [http.ServeMux] used by [Service] with
// all handle functions automatically registered.
var ServiceMux = &serviceMux

var serviceMux http.ServeMux

func init() {
	// Request handlers for the default [serviceMux].
	serviceMux.HandleFunc("/healthz", HealthHandler)
}
