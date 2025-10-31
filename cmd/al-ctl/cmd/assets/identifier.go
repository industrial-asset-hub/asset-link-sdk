/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var credentialFilePath string = ""
var identifierOutputFile string = ""

var IdentifierCmd = &cobra.Command{
	Use:   "identifier",
	Short: "Trigger get identifiers",
	Long:  `This command triggers a get identifiers job and retrieves the identifiers.`,
	Run: func(cmd *cobra.Command, args []string) {
		identifiers, err := al.GetIdentifiers(shared.AssetLinkEndpoint, credentialFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error during get identifiers")
		}

		err = al.WriteIdentifierResponsesFile(identifierOutputFile, identifiers)
		if err != nil {
			log.Fatal().Err(err).Msg("Error writing identifiers to file")
		}
	},
}

func init() {
	IdentifierCmd.Flags().StringVarP(&credentialFilePath, "credential-file", "p", "", "Path to the credential file")
	IdentifierCmd.Flags().StringVarP(&identifierOutputFile, "output", "o", "", "Output filename (default stdout)")
}
