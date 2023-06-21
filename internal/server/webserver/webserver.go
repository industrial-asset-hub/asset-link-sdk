/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package webserver

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "fmt"
  "net/http"
  "os"

  "github.com/gin-gonic/gin"
  "github.com/rs/zerolog/log"
)

// Server for the agent
type Server struct {
  router   *gin.Engine
  endpoint string
  version  metadata.Version
}

func NewServer() *Server {
  return NewServerWithParameters(metadata.Version{
    Version: "unknown",
    Commit:  "unknown",
    Date:    "unknown",
  })
}

func NewServerWithParameters(version metadata.Version) *Server {
  r := gin.Default()
  s := &Server{
    router:   r,
    version:  version,
    endpoint: fmt.Sprintf(":%s", getPORT()),
  }
  s.configureRoutes(r)

  return s
}

func (s *Server) Run() {
  log.Printf("Listening on %s", s.endpoint)
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
}

func getPORT() string {
  port := os.Getenv("PORT")
  if port == "" {
    port = "8090"
  }
  return port
}
