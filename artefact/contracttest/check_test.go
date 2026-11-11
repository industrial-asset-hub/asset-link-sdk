// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package contracttest

import (
	"context"
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	common "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

type fatalProbe struct{ testing.TB }

func (fatalProbe) Fatalf(format string, args ...any) { panic(fmt.Sprintf(format, args...)) }
func TestCheckFixtures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	metadata := &artefact.ArtefactMetaData{JobId: "job", ConnectionParameterSet: &common.ConnectionParameterSet{SchemaUri: "urn:fixture", ParameterJson: `{}`}, OciReference: &generated.OciArtefactReference{RegistryUrl: "http://registry:5000", Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}}
	handler := func(_ *artefact.ArtefactMetaData, p artefact.ProgressReporter) (*artefact.Result, error) {
		return artefact.Succeeded(), p.ReportProgress("download", 100)
	}
	result := Check(t, ctx, metadata, handler, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED)
	require.Equal(t, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED, result.Outcome)
	probe := fatalProbe{TB: t}
	require.Panics(t, func() { Check(probe, ctx, nil, handler, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED) })
	require.Panics(t, func() { Check(probe, ctx, metadata, handler, generated.UpdateOutcome_UPDATE_OUTCOME_FAILED) })
	require.Panics(t, func() {
		Check(probe, ctx, metadata, func(_ *artefact.ArtefactMetaData, p artefact.ProgressReporter) (*artefact.Result, error) {
			return nil, p.ReportProgress("download", 101)
		}, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN)
	})
}
