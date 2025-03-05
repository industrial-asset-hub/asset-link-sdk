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
	var out *os.File = os.Stdout
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
		logger = zerolog.New(out).With().Caller().Logger()
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
	// logLevelRaw := flag.GetString(cli.LogLevel.ToViper())
	// logLevelRaw := flag.Args("log-level")
	lvl, err := zerolog.ParseLevel(logLevelRaw)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid log level format")
	}
	zerolog.SetGlobalLevel(lvl)
}

func SetColorForLogLevel() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	consoleWriter.FormatLevel = func(i interface{}) string {
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				return "\033[34mTRACE\033[0m"
			case "debug":
				return "\033[36mDEBUG\033[0m"
			case "info":
				return "\033[32mINFO\033[0m"
			case "warn":
				return "\033[33mWARN\033[0m"
			case "error":
				return "\033[31mERROR\033[0m"
			case "fatal":
				return "\033[35mFATAL\033[0m"
			case "panic":
				return "\033[31mPANIC\033[0m"
			default:
				return ll
			}
		}
		return ""
	}
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}
