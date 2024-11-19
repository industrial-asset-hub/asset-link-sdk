/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package dcd

import (
	"fmt"

	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/shared"
	driverinfo "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/conn_suite_drv_info"
	discovery "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func PrintInfo(endpoint string) {
	log.Trace().Str("Endpoint", endpoint).Msg("Fetching health")

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	printDriverInfo(conn)
	printDiscoveryOptions(conn)
}

func printDriverInfo(conn *grpc.ClientConn) {
	client := driverinfo.NewDriverInfoApiClient(conn)

	resp, err := client.GetVersionInfo(context.Background(), &driverinfo.GetVersionInfoRequest{})

	if err != nil {
		log.Err(err).Msg("version request returned an error")
		return
	}

	var version = resp.GetVersion()

	var versionNumber = fmt.Sprintf("%d.%d.%d", version.GetMajor(), version.GetMinor(), version.GetPatch())

	println("Asset Link:")
	printVersionInfo("ProductName", version.GetProductName())
	printVersionInfo("ProdcutDescription", version.GetProductDescription())
	printVersionInfo("Version", versionNumber)
	printVersionInfo("VendorName", version.GetVendorName())
	printVersionInfo("Suffix", version.GetSuffix())
	printVersionInfo("DocuURL", version.GetDocuUrl())
	printVersionInfo("FeedbackURL", version.GetFeedbackUrl())
}

func printVersionInfo(key string, value string) {
	if value != "" {
		println("  " + key + ": " + value)
	}
}

func printDiscoveryOptions(conn *grpc.ClientConn) {
	client := discovery.NewDeviceDiscoverApiClient(conn)

	foResp, foErr := client.GetFilterOptions(context.Background(), &discovery.FilterOptionsRequest{})
	if foErr != nil {
		log.Err(foErr).Msg("filter options request returned error")
		return
	}

	ftResp, ftErr := client.GetFilterTypes(context.Background(), &discovery.FilterTypesRequest{})
	if ftErr != nil {
		log.Err(ftErr).Msg("filter types request returned error")
		return
	}

	fos := foResp.GetFilterOptions()
	if len(fos) > 0 {
		println("Filter Options:")
		for _, fo := range fos {
			println("  " + fo.GetKey() + " (" + fo.GetDatatype().String() + ")")
		}
	}

	fts := ftResp.GetFilterTypes()
	if len(fts) > 0 {
		println("Filter Types:")
		for _, ft := range fts {
			println("  " + ft.GetKey() + " (" + ft.GetDatatype().String() + ")")
		}
	}
}
