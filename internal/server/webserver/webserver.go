/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package webserver

import (
	"net/http"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/observability"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
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

	discovery := r.Group("discovery")
	discovery.GET("count", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"jobs": observability.GlobalEvents().GetDiscoveryJobsCount(),
		})
	})
	discovery.GET("started", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"jobs": observability.GlobalEvents().GetDiscoveryJobs(),
		})
	})

}
