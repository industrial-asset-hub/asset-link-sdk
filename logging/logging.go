/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

func SetupLogging() {
	var out = os.Stdout
	var format = "auto"

	var logger zerolog.Logger
	if format == "auto" {
		if term.IsTerminal(int(out.Fd())) {
			format = "pretty"
		} else {
			format = "json"
		}
	}
	switch format {
	case "json":
		logger = zerolog.New(out)
	case "pretty":
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.Stamp,
		})
	default:
		fmt.Fprintf(os.Stderr, "Invalid log format: %s\n", format)
	}
	log.Logger = logger.With().Timestamp().Caller().Logger()
}

func AdjustLogLevel(logLevelRaw string) {
	lvl, err := zerolog.ParseLevel(logLevelRaw)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid log level format")
	}
	zerolog.SetGlobalLevel(lvl)
}
