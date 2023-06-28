/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package webserver

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/observability"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "github.com/gin-gonic/gin"
  "github.com/rs/zerolog/log"
  "net/http"
)

// Server for the agent
type Server struct {
  router   *gin.Engine
  endpoint string
  version  metadata.Version
}

func NewServerWithParameters(address string, version metadata.Version) *Server {
  r := gin.Default()
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
