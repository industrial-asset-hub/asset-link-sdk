/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package webserver

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "encoding/json"
  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

func TestServer(t *testing.T) {
  s := NewServerWithParameters("localhost:8080", metadata.Version{
    Version: "v",
    Commit:  "c",
    Date:    "d",
  })
  t.Run("health endpoint should be available", func(t *testing.T) {
    r := request(t, s.router, "GET", "/health", "")
    assert.Equal(t, http.StatusOK, r.Code)
  })
}

func TestServerVersion(t *testing.T) {
  s := NewServerWithParameters("localhost:8080", metadata.Version{
    Version: "v",
    Commit:  "c",
    Date:    "d",
  })
  r := request(t, s.router, "GET", "/version", "")
  assert.Equal(t, http.StatusOK, r.Code)
  res := map[string]string{}
  assert.NoError(t, json.Unmarshal(r.Body.Bytes(), &res))
  assert.Equal(t, res["version"], "v")
  assert.Equal(t, res["commit"], "c")
  assert.Equal(t, res["date"], "d")
}

func request(t *testing.T, router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
  w := httptest.NewRecorder()
  req, err := http.NewRequest(method, path, strings.NewReader(body))
  require.NoError(t, err)
  router.ServeHTTP(w, req)
  return w
}