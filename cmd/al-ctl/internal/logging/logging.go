/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func SetupLoggingCli(logLevelRaw string) {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:          os.Stderr,
		PartsExclude: []string{"time"},
	})

	lvl, err := zerolog.ParseLevel(logLevelRaw)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid log level format")
	}

	zerolog.SetGlobalLevel(lvl)

	if lvl <= zerolog.DebugLevel {
		log.Logger = logger.With().Caller().Logger()
	} else {
		log.Logger = logger.With().Logger()
	}
}
