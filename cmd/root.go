// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/moderncode-source/vex-svc/vex"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/spf13/cobra"
)

var (
	rootVerbose bool
	rootAddr    string
	rootLevel   uint
)

// Log level determines what source to display logs from.
// Available values are: 0 - silent, 1 - service, 2 - CLI, 3 - service and CLI.
// Controlled by the command-line argument "level".
const (
	logLevelSilent = uint(iota)
	logLevelServiceOnly
	logLevelCmdOnly
	logLevelAll
)

func init() {
	rootCmd.Flags().BoolVarP(&rootVerbose, "verbose", "v", false, "run Vex CLI and service with debug logs")
	rootCmd.Flags().StringVar(&rootAddr, "addr", ":8080", "the TCP network address for the Vex server to listen on")
	rootCmd.Flags().UintVarP(&rootLevel, "level", "l", 3, "log level. 0 - silent, 1 - service only, 2 - CLI only, 3 - both")
}

var rootCmd = &cobra.Command{
	Use:   "vex",
	Short: "Run arbitrary code under isolated environments",
	Long: `Vex is a virtual execution micro-service that runs arbitrary code
in the cloud under controlled, isolated environments.
Documentation is available at https://github.com/moderncode-source/vex-svc`,
	SilenceUsage: true, // Do not print usage on error.
	RunE: func(_ *cobra.Command, _ []string) error {
		// Prepare loggers for this command and the service.
		cmdLogger := globalCmdNopLogger
		svcLogger := globalCmdNopLogger

		// Service uses a thread-safe, fast logger.
		w := diode.NewWriter(logOutput, diodeWriterSize, 0, func(missed int) {
			cmdLogger.Warn().Msgf("Service dropped %d logs", missed)
		})

		// https://github.com/rs/zerolog/tree/master#leveled-logging
		switch rootLevel {
		case logLevelSilent:
		case logLevelServiceOnly:
			if rootVerbose {
				svcLogger = zerolog.New(w)
			} else {
				svcLogger = zerolog.New(w).Level(zerolog.InfoLevel)
			}
		case logLevelCmdOnly:
			if rootVerbose {
				cmdLogger = globalCmdLogger
			} else {
				cmdLogger = globalCmdLogger.Level(zerolog.InfoLevel)
			}
		case logLevelAll:
			if rootVerbose {
				svcLogger = zerolog.New(w)
				cmdLogger = globalCmdLogger
			} else {
				svcLogger = zerolog.New(w).Level(zerolog.InfoLevel)
				cmdLogger = globalCmdLogger.Level(zerolog.InfoLevel)
			}
		default:
			return fmt.Errorf("invalid argument \"%d\" for \"level\" flag", rootLevel)
		}

		cmdLogger.Info().Msg("Welcome to Vex - a virtual execution micro-service")

		// Create a new Vex service.
		var svc *vex.Service
		var err error
		addr := rootAddr

		if !rootVerbose {
			// Avoid the overhead of chaining HTTP request
			// handlers because of the logging middleware.
			svc, err = vex.New(addr, &svcLogger)
		} else {
			// Because [vex.ServiceMux] is a pointer, we can do something like
			// this below, where we register request handlers after we pass
			// the final handler to the service.
			// TODO: put this into a self-contained example.
			withMiddleware := loggingHTTPMiddleware(&svcLogger)
			handler := withMiddleware(vex.ServiceMux)

			svc = vex.NewWithHandler(addr, handler, &svcLogger)
			err = svc.RegisterDefaultHandlers(vex.ServiceMux)
		}

		if err != nil {
			cmdLogger.Error().Err(err).Msg("Service creation error")
			return fmt.Errorf("service creation error: %v", err)
		}

		if err := svc.Validate(); err != nil {
			cmdLogger.Error().Err(err).Msg("Service error")
			return fmt.Errorf("service error: %v", err)
		}

		cmdLogger.Info().Msgf("Starting service process [%d] (Press CTRL+C to quit)", os.Getpid())

		// If an interrupt signal is caught, gracefully shut down the service
		// and encapsulate any error it returns.
		shutdownErr := make(chan error, 1)
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

		var wg sync.WaitGroup
		defer wg.Wait() // Wait for the go-routine below to exit.
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()

			<-ctx.Done()
			cmdLogger.Info().Msg("Shutting down...")

			// TODO: pass a context with timeout.
			if err := svc.Stop(context.Background()); err != nil {
				shutdownErr <- err
			} else {
				close(shutdownErr)
			}
		}()

		if err := svc.Start(); err != nil {
			cmdLogger.Error().Err(err).Msg("Service start error")
			return fmt.Errorf("service start error: %v", err)
		}

		if err := <-shutdownErr; err != nil {
			cmdLogger.Error().Err(err).Msg("Service shutdown error")
			return fmt.Errorf("service shutdown error: %v", err)
		}

		cmdLogger.Info().Msg("Service shutdown complete")
		cmdLogger.Info().Msgf("Finished service process [%d]", os.Getpid())
		return nil
	},
}

// Execute matches ands runs the appropriate CLI command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("vex CLI exited with error: %s", err)
	}
	return nil
}
