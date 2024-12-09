/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package registryclient

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"
)

type appTypes int

const (
	CDM_AGENT                   appTypes = 0
	CDM_DEVICE_CLASS_DRIVER     appTypes = 1
	APP_TYPE_CS_IAH_DISCOVER_V1          = "siemens.industrialassethub.discover.v1"
)

func (apptypes appTypes) String() string {
	return []string{"cdm-agent", "cdm-device-class-driver"}[apptypes]
}

type GrpcServerRegistry struct {
	alId                      string
	grpcServerRegistryAddress string
	grpcAddress               string
	appInstanceId             string
	connection                *grpc.ClientConn
	client                    pb.RegistryApiClient
}

const (
	retryRegistrationInterval            = 10 // Interval in case of errors during an registration [s]
	reRegistrationRefreshInterval uint32 = 60 // Interval when the registration is refreshed [s]
)

// Create new GRPC Registry client
func New(registryAddress string, alId string, grpcAddress string) *GrpcServerRegistry {
	return &GrpcServerRegistry{grpcServerRegistryAddress: registryAddress,
		alId:        alId,
		grpcAddress: grpcAddress,
	}
}

func (r *GrpcServerRegistry) connect() error {
	connection, err := grpc.Dial(r.grpcServerRegistryAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	r.connection = connection
	if err != nil {
		log.Fatal().Err(err).Msg("Could not dial grpc registry")
		return err
	}
	r.client = pb.NewRegistryApiClient(r.connection)
	return nil
}
func (r *GrpcServerRegistry) Stop() {
	log.Debug().
		Msg("Stop gRPC Registry Client")

	r.disconnect()
}

// disconnect closes the connection to the gRPC registry
func (r *GrpcServerRegistry) disconnect() {
	log.Debug().
		Msg("Disconnect from gRPC server registry")

	if r.connection != nil {
		if err := r.unregister(); err != nil {
			log.Fatal().Err(err).Msg("Error during un registering of the service.")
		}
		r.connection.Close()
		r.connection = nil
	}
	if r.client != nil {
		r.client = nil
	}
}

// Register registers the asset link (AL) at the gRPC Server Registry
// It takes also care, to re-new the registration
func (r *GrpcServerRegistry) Register() {
	log.Info().Str("gRPC Server Registry", r.grpcServerRegistryAddress).Msg("Register asset link at grpc server registry")

	// Start registration async. The goroutine also deals with
	// the re-registration at a given interval.
	go func() {
		defer func() {
			if err := r.unregister(); err != nil {
				log.Fatal().Err(err).Msg("Error during unregistering of the service.")
			}
		}()
		for {
			// Try to register, if an errors occur
			retryInterval := reRegistrationRefreshInterval
			err, expireTime := r.register()
			if err != nil {
				if expireTime < reRegistrationRefreshInterval {
					retryInterval = expireTime / 2
				}
			}
			log.Info().
				Uint32("Registration expire time [s]", expireTime).
				Uint32("Re-new registration in [s]", retryInterval).
				Msg("Wait until renew server registration")

			time.Sleep(time.Duration(retryInterval) * time.Second)
		}
	}()
}

func (r *GrpcServerRegistry) unregister() error {
	log.Info().
		Str("App Instance Id", r.appInstanceId).
		Msg("Unregister service")

	unRegisterRequest := pb.UnregisterServiceRequest{Info: &pb.ServiceInfo{
		AppInstanceId: r.appInstanceId,
	}}

	log.Debug().
		Str("Service Info", unRegisterRequest.String()).
		Msg("Unregister service")

	_, err := r.client.UnregisterService(context.Background(), &unRegisterRequest)
	if err != nil {
		log.Warn().Err(err).Msg("An error occured during unregistering of the service.")
		return err
	}
	return nil
}

// register registers the Asset Link at the gprc server registry
func (r *GrpcServerRegistry) register() (error, uint32) {
	if err := r.connect(); err != nil {
		log.Warn().Err(err).Msg("Could not dial GRPC server registry")
		return err, retryRegistrationInterval
	}

	// Split into host and port
	hostName, portNumberString, err := net.SplitHostPort(r.grpcAddress)

	// Catch if no host part is given e.g. *:8080
	if hostName == "" {
		log.Fatal().
			Str("IP or Hostname", hostName).
			Msg("No valid Hostname given. Should be an IP or DNS name.")
	}
	portNumber, _ := strconv.Atoi(portNumberString)
	if err != nil {
		log.Warn().Err(err).Msg("Could parse GRPC server address")
		return err, retryRegistrationInterval
	}

	r.appInstanceId = CDM_DEVICE_CLASS_DRIVER.String() + "-" + r.alId
	register := pb.RegisterServiceRequest{Info: &pb.ServiceInfo{
		AppTypes:         []string{APP_TYPE_CS_IAH_DISCOVER_V1},
		AppInstanceId:    r.appInstanceId,
		DriverSchemaUris: []string{r.alId},
		GrpcIpPortNumber: uint32(portNumber),
	}}
	// Check if IP or DNS name
	if r := net.ParseIP(hostName); r == nil {
		log.Debug().
			Str("Given IP/Hostname", hostName).
			Uint32("Port", uint32(portNumber)).
			Msg("Hostname seems to be an DNS name.")
		register.Info.GrpcIp = &pb.ServiceInfo_DnsDomainname{DnsDomainname: hostName}
	} else {
		log.Debug().
			Str("Given IP/Hostname", hostName).
			Uint32("Port", uint32(portNumber)).
			Msg("Hostname seems to be an IP address name.")
		register.Info.GrpcIp = &pb.ServiceInfo_Ipv4Address{Ipv4Address: hostName}
	}

	log.Debug().
		Str("Service Info", register.String()).
		Msg("Registering service")

	response, err := r.client.RegisterService(context.Background(), &register)
	if err != nil {
		log.Warn().Err(err).Msg("Could not register at grpc server registry")

		return err, retryRegistrationInterval
	} else {
		expireTime := response.ExpireTime
		return nil, expireTime
	}
}
