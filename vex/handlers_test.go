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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/moderncode-source/vex-svc/vex"
)

const mockURL = "http://localhost:8080"

// mockService is used to test request handler responses.
// The actual HTTP server never starts.
var mockService = vex.NewWithHandler("", nil, nil)

func checkHandlerResponseCode(
	ctx context.Context, handler http.Handler, want int, method, url string,
) error {
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(ctx, method, url, nil)

	handler.ServeHTTP(res, req)
	if res.Code != want {
		return fmt.Errorf("Expected response code %d, got %d", want, res.Code)
	}

	return nil
}

func TestHealthAndReadyHandlers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tests := []struct {
		ctx         context.Context
		handler     http.Handler
		want        int
		method, url string
	}{
		{ctx, http.HandlerFunc(mockService.HealthHandler), http.StatusOK, http.MethodGet, mockURL + vex.HealthEndpoint},

		{ctx, http.HandlerFunc(mockService.ReadyHandler), http.StatusOK, http.MethodGet, mockURL + vex.ReadyEndpoint},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s %s", tt.method, tt.url)

		t.Run(testname, func(t *testing.T) {
			if err := checkHandlerResponseCode(
				tt.ctx, tt.handler, tt.want, tt.method, tt.url,
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAndPostQueueHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	url := mockURL + vex.QueueEndpoint

	{
		t.Logf("Testing request to POST %s (empty body)", url)
		res := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")

		mockService.PostQueueHandler(res, req)
		if want := http.StatusBadRequest; res.Code != want {
			t.Fatalf("Expected response code %d, got %d", want, res.Code)
		}
	}

	{
		t.Logf("Testing request to POST %s (invalid body)", url)

		buff := bytes.NewBufferString("{}")

		res := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, url, buff)
		req.Header.Set("Content-Type", "application/json")

		mockService.PostQueueHandler(res, req)
		if want := http.StatusOK; res.Code != want {
			t.Fatalf("Expected response code %d, got %d", want, res.Code)
		}
	}

	{
		t.Logf("Testing request to POST %s (invalid content-type)", url)

		buff := bytes.NewBufferString("{}")

		res := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, url, buff)
		req.Header.Set("Content-Type", "text/plain")

		mockService.PostQueueHandler(res, req)
		if want := http.StatusUnsupportedMediaType; res.Code != want {
			t.Fatalf("Expected response code %d, got %d", want, res.Code)
		}
	}

	t.Logf("Testing requests to GET&POST %s (valid body)...", url)

	for i := 0; ; i++ {
		// End the loop if the queue got full. We are not testing
		// submission processing here, just the handler response.
		res := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, url, nil)

		mockService.GetQueueHandler(res, req)
		if want := http.StatusOK; res.Code != want {
			t.Fatalf("Expected response code %d, got %d", want, res.Code)
		}

		queueResp := struct {
			Length int  `json:"length"`
			Full   bool `json:"full"`
		}{}

		d := json.NewDecoder(res.Body)
		if err := d.Decode(&queueResp); err != nil {
			t.Fatalf("Failed to decode response, err: %s", err)
		}

		if queueResp.Full {
			break
		}

		// Post a valid submission.
		w := &bytes.Buffer{}

		submission := vex.Submission{
			ID:        int64(i + 1),
			Timestamp: vex.TimeUnix{time.Now()},
		}

		e := json.NewEncoder(w)
		if err := e.Encode(submission); err != nil {
			t.Fatalf("Failed to encode submission, err: %s", err)
		}

		res = httptest.NewRecorder()
		req = httptest.NewRequestWithContext(ctx, http.MethodPost, url, w)

		mockService.PostQueueHandler(res, req)
		if want := http.StatusOK; res.Code != want {
			t.Fatalf("Expected response code %d, got %d", want, res.Code)
		}
	}

	{
		// Post the last valid submission while the queue is full (hopefully).
		// We should get informed that the queue is indeed full.

		t.Logf("Testing request to POST %s (valid body, q is full)", url)
		w := &bytes.Buffer{}

		submission := vex.Submission{
			ID:        int64(-1),
			Timestamp: vex.TimeUnix{time.Now()},
		}

		e := json.NewEncoder(w)
		if err := e.Encode(submission); err != nil {
			t.Fatalf("Failed to encode submission, err: %s", err)
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, url, w)

		mockService.PostQueueHandler(res, req)
		if want := http.StatusInsufficientStorage; res.Code != want {
			t.Fatalf("Expected response code %d, got %d", want, res.Code)
		}
	}
}
