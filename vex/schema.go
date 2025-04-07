// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex submission primitives and request/response payload schemas.

package vex

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

type httpReply struct {
	Message string `json:"message"`
	Status  int    `json:"status,string"`
}

// TimeUnix represents an instant in time. It wraps [time.Time], and
// marshals into and unmarshals from a Unix time encoding of a JSON value.
type TimeUnix struct {
	time.Time
}

// MarshalJSON marshals a timestamp in Unix time format into valid JSON.
func (t TimeUnix) MarshalJSON() ([]byte, error) {
	sec := t.Unix()
	s := strconv.FormatInt(sec, 10) // Base 10.
	return unsafe.Slice(unsafe.StringData(s), len(s)), nil
}

// UnmarshalJSON unmarshals the valid Unix time encoding of a JSON value b.
func (t *TimeUnix) UnmarshalJSON(b []byte) error {
	// Avoid making a copy of the given byte slice that would arise from
	// casting the byte slice to string type.
	s := unsafe.String(&b[0], len(b))

	sec, err := strconv.ParseInt(s, 10, 64) // Base 10, 64 bits.
	if err != nil {
		return fmt.Errorf("TimeUnix.UnmarshalJSON error: %s", err)
	}

	t.Time = time.Unix(sec, 0)
	return nil
}

// Submission type represents a code execution request data sent by the client.
// Submissions are queued by the service and later scheduled to execute.
type Submission struct {
	// ID is a unique identifier for this submission.
	ID int64 `json:"id"`

	// Timestamp stores the time in Unix format at which this submission was
	// created by the client. Used to monitor and estimate an average total
	// time it takes to execute a submission with respect to current traffic.
	Timestamp TimeUnix `json:"timestamp"`

	// TODO: add more fields and validate them in Validate().
}

// Validate reports whether the submission's fields values are valid.
func (s *Submission) Validate() bool {
	return true
}
