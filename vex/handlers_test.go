// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex HTTP request handlers tests.

package vex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moderncode-source/vex-svc/vex"
)

func TestHealthHandler(t *testing.T) {
	url := "http://localhost:8080/healthz"
	req := httptest.NewRequestWithContext(context.Background(), "", url, nil)
	res := httptest.NewRecorder()

	vex.HealthHandler(res, req)

	want := 200

	req.Method = http.MethodGet
	if res.Code != want {
		t.Fatalf("Expected response code %d, got %d", want, res.Code)
	}

	req.Method = http.MethodPost
	if res.Code != want {
		t.Fatalf("Expected response code %d, got %d", want, res.Code)
	}
}
