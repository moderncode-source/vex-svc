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
	"net/http"

	"github.com/goccy/go-json"
)

var okReply []byte

func init() {
	// Pre-encode a short "ok" HTTP response body for
	// request handlers that may need it often.
	b, err := json.Marshal(httpReply{Message: "ok", Status: http.StatusOK})
	if err != nil {
		panic(err)
	}
	okReply = b
}

// TODO: add a Strict-Transport-Security headers to every handler.
// TODO (MAIN-107): reuse JSON buffers using [sync.Pool].

// HealthHandler handles requests to service liveness probe endpoint that can
// be used to check whether the server is running.
func (svc *Service) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Cache-Control", "no-store")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(okReply); err != nil {
		svc.logger.Err(err).
			Int("status", http.StatusOK).
			Msg("HealthHandler JSON encoder error")
	}
}

// ReadyHandler handles requests to service readiness probe endpoint that can
// be used to check whether the server is ready to receive traffic.
func (svc *Service) ReadyHandler(w http.ResponseWriter, _ *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Cache-Control", "no-store")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(okReply); err != nil {
		svc.logger.Err(err).
			Int("status", http.StatusOK).
			Msg("ReadyHandler JSON encoder error")
	}
}

// PostQueueHandler handles requests that
// post a new item into the submission queue.
func (svc *Service) PostQueueHandler(w http.ResponseWriter, req *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "deny")
	h.Set("Cache-Control", "no-store")

	var reply httpReply

	defer func() {
		w.WriteHeader(reply.Status)

		e := json.NewEncoder(w)
		if err := e.Encode(reply); err != nil {
			svc.logger.Err(err).
				Int("status", reply.Status).
				Msg("PostQueueHandler JSON encoder error")
		}
	}()

	// Validate client's body content-type.
	// If it is not set (empty), assume "application/json".
	switch req.Header.Get("Content-Type") {
	case "", "application/json", "application/json; charset=utf-8":
	default:
		reply = httpReply{
			Message: "Unsupported 'Content-Type' header. Supported is 'application/json'",
			Status:  http.StatusUnsupportedMediaType,
		}
		return
	}

	// Decode request body into a [Submission].
	var submission Submission
	d := json.NewDecoder(req.Body)

	if err := d.Decode(&submission); err != nil || !submission.Validate() {
		// Failed to decode due to a bad request body
		// or submission is in invalid state. Reply to the client.
		reply = httpReply{
			Message: "Bad request body",
			Status:  http.StatusBadRequest,
		}
		return
	}

	// Insert submission into the queue if it is not full.
	// Inform the client otherwise.
	select {
	case svc.queue <- submission:
		reply = httpReply{
			Message: "Submission successfully posted",
			Status:  http.StatusOK,
		}
	default:
		reply = httpReply{
			Message: "Submission queue is full. Retry later",
			Status:  http.StatusInsufficientStorage,
		}
	}
}

// GetQueueHandler handles requests to the submission queue
// endpoint to retrieve information about the queue.
func (svc *Service) GetQueueHandler(w http.ResponseWriter, _ *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "deny")
	h.Set("Cache-Control", "no-store")

	// Respond with the total number of submissions in the queue
	// and whether the queue is full.

	length := len(svc.queue)
	reply := struct {
		Length int  `json:"length"`
		Full   bool `json:"full"`
	}{Length: length, Full: length == cap(svc.queue)}

	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	if err := e.Encode(reply); err != nil {
		svc.logger.Err(err).
			Int("status", http.StatusOK).
			Msg("GetQueueHandler JSON encoder error")
	}
}
