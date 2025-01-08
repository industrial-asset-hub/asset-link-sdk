/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
)

type ArtefactIdentifier struct {
	artefactType *generated.ArtefactType
}

func NewArtefactIdentifier(artefactType *generated.ArtefactType) *ArtefactIdentifier {
	artefactIdentifier := &ArtefactIdentifier{artefactType: artefactType}
	return artefactIdentifier
}

func (ai *ArtefactIdentifier) GetArtefactType() *generated.ArtefactType {
	return ai.artefactType
}
