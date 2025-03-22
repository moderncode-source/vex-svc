// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex service HTTP request handlers.

package vex

import (
	"encoding/json" // TODO: consider using a faster alternative.
	"fmt"
	"net/http"
)

// HealthHandler handles requests to service liveness probe endpoint that can
// be used to check whether the server is running.
func (svc *Service) HealthHandler(w http.ResponseWriter, req *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Cache-Control", "no-store")

	// TODO: add a Strict-Transport-Security header.

	if len(req.Method) != 0 && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ReadyHandler handles requests to service readiness probe endpoint that can
// be used to check whether the server is ready to receive traffic.
func (svc *Service) ReadyHandler(w http.ResponseWriter, req *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Cache-Control", "no-store")

	// TODO: add a Strict-Transport-Security header.

	if len(req.Method) != 0 && req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PostQueueHandler handles requests that
// post a new item into the submission queue.
func (svc *Service) PostQueueHandler(w http.ResponseWriter, req *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "deny")
	h.Set("X-Xss-Protection", "1; mode=block")
	h.Set("Cache-Control", "no-store")
	h.Set("Location", QueueEndpoint)

	// TODO: add a Strict-Transport-Security header.

	// TODO: respect Accept header value.
	//       Maybe they only want to receive a json response.

	// TODO: be more lenient here. Assume json if content-type is not set,
	//       accept "application/json; charset=utf-8", and try to decode
	//       "text/plain" into [Submission] too.
	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// Decode request body into a [Submission].
	var submission Submission
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&submission); err != nil {
		// TODO: consider responding with the decoding error message here.
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: replace queue with a proper construct.
	queue = append(queue, submission)
	w.WriteHeader(http.StatusOK)
}

// GetQueueHandler handles requests to the submission queue
// endpoint to retrieve information about the queue.
func (svc *Service) GetQueueHandler(w http.ResponseWriter, req *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "deny")
	h.Set("X-Xss-Protection", "1; mode=block")
	h.Set("Cache-Control", "no-store")

	// TODO: add a Strict-Transport-Security header.

	// TODO: respect Accept header value.
	//       Maybe they only want to receive a json response.

	// Respond with the total number of submissions in the queue.
	// TODO: respond with an array of submission ids in the queue instead.
	if n, err := fmt.Fprint(w, len(queue)); err != nil && n == 0 {
		h.Del("Content-Length")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    // TODO: write the status code before writing body into w.
    //       ignore the error above, only log it.
	w.WriteHeader(http.StatusOK)
}
