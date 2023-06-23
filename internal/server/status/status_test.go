package status

import (
  generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/status"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "github.com/stretchr/testify/assert"
  "golang.org/x/net/context"
  "testing"
)

func TestStatusHealth(t *testing.T) {
  s := StatusServerEntity{}

  resp, err := s.GetHealth(context.Background(), &generated.HealthRequest{})
  assert.NoError(t, err)
  assert.Equal(t, resp, &generated.HealthReply{Health: generated.HealthReply_HEALTHY})
}

func TestStatusVersion(t *testing.T) {
  s := StatusServerEntity{Version: metadata.Version{
    Version: "1",
    Commit:  "commit",
    Date:    "date",
  }}

  resp, err := s.GetVersion(context.Background(), &generated.VersionRequest{})
  assert.NoError(t, err)
  assert.Equal(t, resp, &generated.VersionReply{Version: "1", Commit: "commit", Date: "date"})
}
