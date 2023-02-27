/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 * Author: {{cookiecutter.author_name}} <{{cookiecutter.author_email}}>
 */

package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

func SetupLogging() {
	var writer io.Writer = os.Stdout
	if term.IsTerminal(int(os.Stdout.Fd())) {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Stamp,
		}
	}
	log.Logger = zerolog.New(writer).With().
		Timestamp().
		Caller().
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
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
