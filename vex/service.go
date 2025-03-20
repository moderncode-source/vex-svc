// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex service.

package vex

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Vex major, minor, and patch version numbers.
const (
	VersionMajor = 0
	VersionMinor = 0
	VersionPatch = 1
)

// Block clients from keeping connections open forever by setting
// a deadline for reading request headers. See [http.Server].
const serverReadHeaderTimeout = 10 * time.Second

type Service struct {
	server *http.Server
}

// New allocates and returns a new [Service] with [http.Server] that will
// listen on TCP network address addr and handle requests on incoming
// connections using [ServiceMux] handler. This is the recommended and default
// way to create a Vex service.
//
// To choose your own handler or fall back to [http.DefaultServeMux],
// use [NewWithHandler].
func New(addr string) *Service {
	return &Service{
		server: &http.Server{
			ReadHeaderTimeout: serverReadHeaderTimeout,
			Addr:              addr,
			Handler:           ServiceMux,
		},
	}
}

// NewWithHandler allocates and returns a new [Service] with [http.Server]
// that will listen on TCP network address addr and handle requests on
// incoming connections by calling [http.Server.Serve] with handler.
//
// If handler is nil, [http.DefaultServeMux] will be used.
//
// See also: [New].
func NewWithHandler(addr string, handler http.Handler) *Service {
	return &Service{
		server: &http.Server{
			ReadHeaderTimeout: serverReadHeaderTimeout,
			Addr:              addr,
			Handler:           handler,
		},
	}
}

func (svc *Service) Start() error {
	err := svc.server.ListenAndServe()
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("failed to serve service: %v", err)
}

func (svc *Service) Stop(ctx context.Context) error {
	err := svc.server.Shutdown(ctx)
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("failed to stop service: %v", err)
}
