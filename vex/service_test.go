// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex service tests.

package vex_test

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/moderncode-source/vex-svc/vex"
)

func TestNew(t *testing.T) {
	// TODO: use httptest.NewRequestWithContext().

	svc := vex.New(":8080")
	var wg sync.WaitGroup

	// Give the entire operation with the service to complete within 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	defer func() {
		if err := svc.Stop(ctx); err != nil {
			t.Logf("Failed to stop service: %s", err)
		}
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		if err := svc.Start(); err != nil {
			t.Logf("Server exited with error: %s", err)
		}
		wg.Done()
	}()

	url := "http://localhost:8080/healthz"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		t.Fatalf("Request returned error: %s", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response code 200, got: %d", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %s", err)
		}
	}()
}
