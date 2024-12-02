/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package webserver

import (
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/observability"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Server for the agent
type Server struct {
	router   *gin.Engine
	endpoint string
	version  metadata.Version
}

func NewServerWithParameters(address string, version metadata.Version) *Server {
	r := gin.New()
	r.Use(logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
		return log.Logger
	})))
	s := &Server{
		router:   r,
		version:  version,
		endpoint: address,
	}
	s.configureRoutes(r)

	return s
}

func (s *Server) Run() {
	if err := s.router.Run(s.endpoint); err != nil {
		log.Fatal().Err(err).Msg("could not start webserver")
	}
}

func (s *Server) configureRoutes(r *gin.Engine) {
	r.GET("health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("version", func(context *gin.Context) {
		context.JSON(http.StatusOK, s.version)
	})

	r.GET("stats", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"discovery_jobs":         observability.GlobalEvents().GetDiscoveryJobsCount(),
			"discovery_jobs_started": observability.GlobalEvents().GetDiscoveryJobs(),
		})
	})
}
