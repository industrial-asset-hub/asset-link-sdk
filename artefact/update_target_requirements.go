/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefact

import (
	"context"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
)

// UpdateTargetRequirementsProvider is an optional capability of an Update
// implementation. Its response describes how the cloud can construct a device
// identifier from the selected asset's Inventory v2 read representation.
// It must not return device-specific values in this response.
type UpdateTargetRequirementsProvider interface {
	GetUpdateTargetRequirements(context.Context, *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error)
}
