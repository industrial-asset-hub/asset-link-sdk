/*******************************************************************************
* Copyright (c) Siemens AG 2022 ALL RIGHTS RESERVED.
*******************************************************************************/

package server

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	pb "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/conn_suite_registry"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/registry/internal/shared"

	"github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type serviceEntry struct {
	driver    *pb.ServiceInfo
	timeAdded time.Time
}

var (
	serviceExpireTime uint32 = 60 * 15 // [s] 15 mins
	// map[string]serviceEntry{})
	registeredServices sync.Map
)

// Server implements the gRPC Server
type Server struct {
	pb.UnimplementedRegistryApiServer
}

func NewServer() pb.RegistryApiServer {
	return &Server{}
}

// RegisterService implements the gRPC server method
func (s *Server) RegisterService(ctx context.Context, in *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	if in == nil {
		errMsg := "Input args must not be nil"
		log.Warn().
			Msg(errMsg)
		return &pb.RegisterServiceResponse{}, errors.New(errMsg)
	}

	if in.Info == nil {
		errMsg := "No info{} value given. Nothing to insert."
		log.Warn().
			Msg(errMsg)
		return &pb.RegisterServiceResponse{}, status.Error(codes.InvalidArgument, errMsg)
	}

	if strings.TrimSpace(in.Info.AppInstanceId) == "" {
		errMsg := "No AppInstanceId value given. Abort."
		log.Warn().
			Msg(errMsg)
		return &pb.RegisterServiceResponse{}, status.Error(codes.InvalidArgument, errMsg)
	}

	// Check if endpoint is set
	if (strings.TrimSpace(in.Info.GetIpv4Address()) == "") && (strings.TrimSpace(in.Info.GetDnsDomainname()) == "") {
		errMsg := "No Endpoint value given. Abort."
		log.Warn().
			Msg(errMsg)
		return &pb.RegisterServiceResponse{}, status.Error(codes.InvalidArgument, errMsg)
	}

	// Check if port has a valid range
	if (in.Info.GrpcIpPortNumber < 1) || (in.Info.GrpcIpPortNumber > 65535) {
		errMsg := "No valid Port given. Abort."
		log.Warn().
			Msg(errMsg)
		return &pb.RegisterServiceResponse{}, status.Error(codes.InvalidArgument, errMsg)
	}

	log.Info().
		Str("AppInstanceIds", in.Info.AppInstanceId).
		Str("App Type", strings.Join(in.Info.AppTypes, ", ")).
		Str("Interfaces", strings.Join(in.Info.Interfaces, ", ")).
		Str("Driver Schema Uri", strings.Join(in.Info.DriverSchemaUris, ", ")).
		Str("IPv4", in.Info.GetIpv4Address()).
		Str("DNS Hostname", in.Info.GetDnsDomainname()).
		Int("Port", int(in.Info.GrpcIpPortNumber)).
		Msg("Register service")

	// Check if the entry already exists
	if service, found := registeredServices.Load(in.Info.AppInstanceId); found {
		if reflect.DeepEqual(service.(serviceEntry).driver, in.Info) {
			// No change, re-registering the same service.
			log.Debug().
				Msg("Service will be updated, with same metadata.")
		} else {
			log.Warn().
				Str("AppInstanceId", in.Info.AppInstanceId).
				Msg("AppInstanceId already existing. But, someone tries to modify the registered service with changed values.")
		}

	} else {
		log.Info().
			Str("AppInstanceId", in.Info.AppInstanceId).
			Msg("Registered AppInstanceId")
	}

	// Set default interfaces if not set
	if in.Info.Interfaces == nil || len(in.Info.Interfaces) == 0 {
		log.Info().
			Msg("No interfaces configured. Set default interfaces for the given app types.")
		for _, appType := range in.Info.AppTypes {
			switch appType {
			case "cs-driver":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.connectivitysuite.driverevent.v1",
					"siemens.connectivitysuite.config.v1",
					"siemens.connectivitysuite.data.v1")
			case "noncs-driver":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.connectivitysuite.config.v1")
			case "cs-gateway":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.connectivitysuite.config.v1")
			case "cs-import-converter":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.connectivitysuite.importconverter.v1")
			case "cs-tagbrowser":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.connectivitysuite.browsing.v2")
			case "cs-alarms-events":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.connectivitysuite.alarmsevents.v1")
			case "iah-discover":
				in.Info.Interfaces = append(in.Info.Interfaces,
					"siemens.connectivitysuite.drvinfo.v1",
					"siemens.industrialassethub.discover.v1")
			}
		}

		// Remove duplicates after appending new interfaces
		in.Info.Interfaces = removeDuplicates(in.Info.Interfaces)
	}

	registeredServices.Store(in.Info.AppInstanceId, serviceEntry{driver: in.Info, timeAdded: time.Now()})
	return &pb.RegisterServiceResponse{ExpireTime: serviceExpireTime}, nil
}

// UnregisterService implements the gRPC server method
func (s *Server) UnregisterService(ctx context.Context, in *pb.UnregisterServiceRequest) (*pb.UnregisterServiceResponse, error) {
	if in == nil {
		errMsg := "Input args must not be nil"
		log.Warn().
			Msg(errMsg)
		return &pb.UnregisterServiceResponse{}, status.Error(codes.Internal, errMsg)
	}

	if in.Info == nil {
		errMsg := "No info{} value given. Nothing to delete."
		log.Warn().
			Msg(errMsg)
		return &pb.UnregisterServiceResponse{}, status.Error(codes.InvalidArgument, errMsg)
	}

	log.Info().
		Str("AppInstanceIds", in.Info.AppInstanceId).
		Msg("Unregister service")

	// Search if the entry already exist
	if _, found := registeredServices.LoadAndDelete(in.Info.AppInstanceId); found {
		log.Debug().
			Str("AppInstanceId", in.Info.AppInstanceId).
			Msg("Delete AppInstanceId")
	} else {
		errMsg := "AppInstanceId not found inside registered services"
		log.Warn().
			Msg(errMsg)
		return &pb.UnregisterServiceResponse{}, status.Error(codes.InvalidArgument, errMsg)
	}

	return &pb.UnregisterServiceResponse{}, nil
}

// QueryRegisteredServices Implementation of the Query Registered Service gRPC method
func (s *Server) QueryRegisteredServices(ctx context.Context, in *pb.QueryRegisteredServicesRequest,
) (*pb.QueryRegisteredServicesResponse, error) {

	if in == nil {
		errMsg := "Input args must not be nil"
		log.Warn().
			Msg("Input args must not be nil.")
		return &pb.QueryRegisteredServicesResponse{}, errors.New(errMsg)
	}

	appType := []string{}
	appInstanceIds := []string{}
	driverSchemaUris := []string{}
	interfaces := []string{}

	// No filter set
	var filterEnabled = false
	if in.Filter != nil {
		appType = in.Filter.AppTypes
		appInstanceIds = in.Filter.AppInstanceIds
		driverSchemaUris = in.Filter.DriverSchemaUris
		interfaces = in.Filter.Interfaces

		filterEnabled = true
	}

	log.Info().
		Str("AppTypes", strings.Join(appType, ", ")).
		Str("AppInstanceIds", strings.Join(appInstanceIds, ", ")).
		Str("DriverSchemas", strings.Join(driverSchemaUris, ", ")).
		Str("Interfaces", strings.Join(interfaces, ", ")).
		Msg("Query registered services")
	filteredServices := []*pb.ServiceInfo{}

	// Iterates over map of registered service, and apply filters
	registeredServices.Range(func(k, v interface{}) bool {
		service := v.(serviceEntry)

		// Check if given service is expired
		if checkExpired(service) {
			log.Info().
				Str("appInstanceId", service.driver.AppInstanceId).
				Msg("Timer expired. Removing service from registry.")
			registeredServices.Delete(v.(serviceEntry).driver.AppInstanceId)
			return true
		}

		// Check if filters are enabled or filter are matching
		// If the entry is valid, true is returned as [...]Found
		appTypeFound := len(appType) == 0 ||
			shared.SearchSliceInSlice(appType, service.driver.AppTypes)
		appInstanceIDFound := len(appInstanceIds) == 0 ||
			shared.SearchStringInSlice(appInstanceIds, service.driver.AppInstanceId)
		driverSchemaUrisFound := len(driverSchemaUris) == 0 ||
			shared.SearchSliceInSlice(driverSchemaUris, service.driver.DriverSchemaUris)
		interfacesFound := len(interfaces) == 0 ||
			shared.SearchSliceInSlice(interfaces, service.driver.Interfaces)

		log.Debug().
			Bool("filterEnabled", filterEnabled).
			Bool("appTypeFound", appTypeFound).
			Bool("appInstanceIDFound", appInstanceIDFound).
			Bool("driverSchemaUrisFound", driverSchemaUrisFound).
			Bool("interfacesFound", interfacesFound).
			Msg("Matched filters.")

		if !filterEnabled || (appTypeFound && appInstanceIDFound && driverSchemaUrisFound && interfacesFound) {
			filteredServices = append(filteredServices, service.driver)
			log.Debug().
				Str("Service ", service.driver.AppInstanceId).
				Msg("Found service with given filter settings. Adding the service to return list.")
		}
		return true
	})

	return &pb.QueryRegisteredServicesResponse{Infos: filteredServices}, nil
}
