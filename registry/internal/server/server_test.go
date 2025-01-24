/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package server

import (
	"context"
	"sort"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"
	"github.com/stretchr/testify/assert"
)

func eraseStorage() {
	registeredServices.Range(func(key interface{}, value interface{}) bool {
		registeredServices.Delete(key)
		return true
	})
}

func initialStorageList() (validService int) {
	// Prepare server registry list

	// Driver 1
	in := &pb.ServiceInfo{AppInstanceId: "driver-1",
		GrpcIp: &pb.ServiceInfo_Ipv4Address{
			Ipv4Address: "1.2.3.1",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("driver-1", serviceEntry{driver: in, timeAdded: time.Now()})

	// Driver 2
	in = &pb.ServiceInfo{AppInstanceId: "driver-2",
		GrpcIp: &pb.ServiceInfo_Ipv4Address{
			Ipv4Address: "1.2.3.2",
		},
		GrpcIpPortNumber: 1234,
		DriverSchemaUris: []string{"uri1"},
	}
	registeredServices.Store("driver-2", serviceEntry{driver: in, timeAdded: time.Now()})

	// Driver 3
	in = &pb.ServiceInfo{AppInstanceId: "driver-3",
		AppTypes: []string{"type1", "type2"},
		GrpcIp: &pb.ServiceInfo_Ipv4Address{
			Ipv4Address: "1.2.3.3",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("driver-3", serviceEntry{driver: in, timeAdded: time.Now()})

	// Driver 4
	in = &pb.ServiceInfo{AppInstanceId: "driver-4",
		AppTypes: []string{"type2"},
		GrpcIp: &pb.ServiceInfo_Ipv4Address{
			Ipv4Address: "1.2.3.4",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("driver-4", serviceEntry{driver: in, timeAdded: time.Now()})

	// Driver 5
	in = &pb.ServiceInfo{AppInstanceId: "driver-5",
		AppTypes:         []string{"type3"},
		DriverSchemaUris: []string{"uri2"},
		GrpcIp: &pb.ServiceInfo_Ipv4Address{
			Ipv4Address: "1.2.3.5",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("driver-5", serviceEntry{driver: in, timeAdded: time.Now()})

	// Adding an expired driver
	in = &pb.ServiceInfo{AppInstanceId: "expired-driver",
		AppTypes: []string{"expired"},
		GrpcIp: &pb.ServiceInfo_DnsDomainname{
			DnsDomainname: "driver.expired.com",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("expired-driver", serviceEntry{driver: in, timeAdded: time.Now().Add(-901 * time.Second)})

	// Adding driver which will exprire in 1 second
	in = &pb.ServiceInfo{AppInstanceId: "not-expired-driver-2",
		AppTypes: []string{"expired"},
		GrpcIp: &pb.ServiceInfo_DnsDomainname{
			DnsDomainname: "driver2.expired.com",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("not-expired-driver-2", serviceEntry{driver: in, timeAdded: time.Now().Add(-899 * time.Second)})

	// Adding a driver with interfaces
	in = &pb.ServiceInfo{AppInstanceId: "driver-8",
		AppTypes:         []string{"type4"},
		Interfaces:       []string{"package1", "package2", "package3"},
		DriverSchemaUris: []string{"uri3"},
		GrpcIp: &pb.ServiceInfo_Ipv4Address{
			Ipv4Address: "1.2.3.8",
		},
		GrpcIpPortNumber: 1234,
	}
	registeredServices.Store("driver-8", serviceEntry{driver: in, timeAdded: time.Now()})

	return 7
}

func registerDriversForAutocomplete() (validService int) {

	s := NewServer()
	ctx := context.Background()

	// Adding drivers with app types but without interfaces, to testsuite autocomplete
	serviceInfo := &pb.ServiceInfo{
		AppInstanceId:    "driver-autocomplete1",
		AppTypes:         []string{"iah-discover"},
		DriverSchemaUris: []string{"uri4"},
		GrpcIp: &pb.ServiceInfo_DnsDomainname{
			DnsDomainname: "driver1.autocomplete.com",
		},
		GrpcIpPortNumber: 1234,
	}

	_, err := s.RegisterService(ctx, &pb.RegisterServiceRequest{Info: serviceInfo})
	if err != nil {
		// Handle error
		return 0
	}

	serviceInfo = &pb.ServiceInfo{
		AppInstanceId:    "driver-autocomplete2",
		AppTypes:         []string{"cs-driver", "iah-discover"},
		DriverSchemaUris: []string{"uri5"},
		GrpcIp: &pb.ServiceInfo_DnsDomainname{
			DnsDomainname: "driver2.autocomplete.com",
		},
		GrpcIpPortNumber: 1234,
	}

	_, err = s.RegisterService(ctx, &pb.RegisterServiceRequest{Info: serviceInfo})
	if err != nil {
		// Handle error
		return 0
	}

	return 2
}

func TestRegisterService(t *testing.T) {

	eraseStorage()

	type expectation struct {
		out *pb.RegisterServiceResponse
		err error
	}
	testCases := []struct {
		name     string
		req      *pb.RegisterServiceRequest
		expected expectation
	}{
		{
			name: "Smoketest - IPv4",
			req: &pb.RegisterServiceRequest{
				Info: &pb.ServiceInfo{
					AppInstanceId: "Instance A",
					GrpcIp: &pb.ServiceInfo_Ipv4Address{
						Ipv4Address: "1.2.3.4",
					},
					GrpcIpPortNumber: 1234,
				}},
			expected: expectation{
				out: &pb.RegisterServiceResponse{ExpireTime: 900},
				err: nil,
			},
		},
		{
			name: "Smoketest - DNS",
			req: &pb.RegisterServiceRequest{
				Info: &pb.ServiceInfo{
					AppInstanceId: "Instance B",
					GrpcIp: &pb.ServiceInfo_DnsDomainname{
						DnsDomainname: "iah.testsuite",
					},
					GrpcIpPortNumber: 1234,
				}},
			expected: expectation{
				out: &pb.RegisterServiceResponse{ExpireTime: 900},
				err: nil,
			},
		},
		{
			name: "NoDNSorIPGiven",
			req: &pb.RegisterServiceRequest{
				Info: &pb.ServiceInfo{
					AppInstanceId:    "Instance C",
					GrpcIpPortNumber: 1234,
				}},
			expected: expectation{
				out: &pb.RegisterServiceResponse{ExpireTime: 900},
				err: status.Error(codes.InvalidArgument, "No Endpoint value given. Abort."),
			},
		},
		{
			name: "NoPortGiven",
			req: &pb.RegisterServiceRequest{
				Info: &pb.ServiceInfo{
					AppInstanceId: "Instance B",
					GrpcIp: &pb.ServiceInfo_DnsDomainname{
						DnsDomainname: "iah.testsuite",
					},
				}},
			expected: expectation{
				out: &pb.RegisterServiceResponse{ExpireTime: 900},
				err: status.Error(codes.InvalidArgument, "No valid Port given. Abort."),
			},
		},
		{
			name: "NoInstanceIdGiven",
			req: &pb.RegisterServiceRequest{
				Info: &pb.ServiceInfo{
					AppInstanceId: "",
				}},
			expected: expectation{
				out: &pb.RegisterServiceResponse{ExpireTime: 900},
				err: status.Error(codes.InvalidArgument, "No AppInstanceId value given. Abort."),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			ctx := context.Background()
			_, err := NewServer().RegisterService(ctx, testCase.req)
			if testCase.expected.err == nil {
				assert.NoError(t, testCase.expected.err)
				_, found := registeredServices.Load(testCase.req.Info.AppInstanceId)
				assert.True(t, found)
			} else {
				assert.Equal(t, testCase.expected.err, err)
			}

		})
	}

}

func TestQueryRegisteredService(t *testing.T) {

	eraseStorage()

	validService := registerDriversForAutocomplete()
	validService += initialStorageList()

	type expectation struct {
		out []*pb.ServiceInfo
		err error
	}

	testCases := []struct {
		name     string
		req      *pb.QueryRegisteredServicesRequest
		expected expectation
	}{
		{

			name: "Smoketest",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "driver-1", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.1"}, GrpcIpPortNumber: 1234},
					{AppInstanceId: "driver-2", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.2"}, GrpcIpPortNumber: 1234, DriverSchemaUris: []string{"uri1"}},
					{AppInstanceId: "driver-3", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.3"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type1", "type2"}},
					{AppInstanceId: "driver-4", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.4"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type2"}},
					{AppInstanceId: "driver-5", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.5"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type3"}, DriverSchemaUris: []string{"uri2"}},
					{AppInstanceId: "driver-8", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.8"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type4"}, Interfaces: []string{"package1", "package2", "package3"}, DriverSchemaUris: []string{"uri3"}},
					{AppInstanceId: "driver-autocomplete1", GrpcIp: &pb.ServiceInfo_DnsDomainname{DnsDomainname: "driver1.autocomplete.com"}, GrpcIpPortNumber: 1234, AppTypes: []string{"iah-discover"}, Interfaces: []string{"siemens.connectivitysuite.drvinfo.v1", "siemens.industrialassethub.discover.v1"}, DriverSchemaUris: []string{"uri4"}},
					{AppInstanceId: "driver-autocomplete2", GrpcIp: &pb.ServiceInfo_DnsDomainname{DnsDomainname: "driver2.autocomplete.com"}, GrpcIpPortNumber: 1234, AppTypes: []string{"cs-driver", "iah-discover"}, Interfaces: []string{"siemens.connectivitysuite.drvinfo.v1", "siemens.connectivitysuite.driverevent.v1", "siemens.connectivitysuite.config.v1", "siemens.connectivitysuite.data.v1", "siemens.industrialassethub.discover.v1"}, DriverSchemaUris: []string{"uri5"}},
					{AppInstanceId: "not-expired-driver-2", GrpcIp: &pb.ServiceInfo_DnsDomainname{DnsDomainname: "driver2.expired.com"}, GrpcIpPortNumber: 1234, AppTypes: []string{"expired"}}},
				err: nil,
			},
		},
		{
			name: "OnlyAppInstanceIdOneSearchEntry",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppInstanceIds: []string{"driver-2"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					&pb.ServiceInfo{AppInstanceId: "driver-2", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.2"}, GrpcIpPortNumber: 1234, DriverSchemaUris: []string{"uri1"}}},
				err: nil,
			},
		},
		{
			name: "OnlyAppInstanceIdMultipleSearchEntries",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppInstanceIds: []string{"driver-2", "driver-3"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "driver-2", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.2"}, GrpcIpPortNumber: 1234, DriverSchemaUris: []string{"uri1"}},
					{AppInstanceId: "driver-3", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.3"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type1", "type2"}},
				},
				err: nil,
			},
		},
		{
			name: "OnlyAppInstanceIdNoMatch",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppInstanceIds: []string{"notfound"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{},
				err: nil,
			},
		},
		{
			name: "OnlyAppTypeFilterOneSearchEntry",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppTypes: []string{"type3"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "driver-5", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.5"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type3"}, DriverSchemaUris: []string{"uri2"}}},
				err: nil,
			},
		},
		{
			name: "OnlyAppTypeFilterMultipleSearchEntries",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppTypes: []string{"type2"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "driver-3", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.3"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type1", "type2"}},
					{AppInstanceId: "driver-4", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.4"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type2"}}},
				err: nil,
			},
		},
		{
			name: "OnlyAppTypeFilterNoMatches",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppTypes: []string{"notfound"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{},
				err: nil,
			},
		},

		{
			name: "OnlyDriverSchemaUriFilterNoMatch",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppInstanceIds: []string{"notfound"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{},
				err: nil,
			},
		},
		{
			name: "OnlyDriverSchemaUriFilterSearchEntry",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{DriverSchemaUris: []string{"uri1"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "driver-2", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.2"}, GrpcIpPortNumber: 1234, DriverSchemaUris: []string{"uri1"}}},
				err: nil,
			},
		},
		{
			name: "OnlyDriverSchemaUriFilterSearchEntries",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{DriverSchemaUris: []string{"uri2"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "driver-5", GrpcIp: &pb.ServiceInfo_Ipv4Address{Ipv4Address: "1.2.3.5"}, GrpcIpPortNumber: 1234, AppTypes: []string{"type3"}, DriverSchemaUris: []string{"uri2"}}},
				err: nil,
			},
		},
		{
			name: "OnlyDriverSchemaUriFilterNoMatches",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{DriverSchemaUris: []string{"notfound"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{},
				err: nil,
			},
		},
		{
			name: "SmokeTestExpiredService",
			req:  &pb.QueryRegisteredServicesRequest{Filter: &pb.ServiceInfoFilter{AppTypes: []string{"expired"}}},
			expected: expectation{
				out: []*pb.ServiceInfo{
					{AppInstanceId: "not-expired-driver-2", GrpcIp: &pb.ServiceInfo_DnsDomainname{DnsDomainname: "driver2.expired.com"}, GrpcIpPortNumber: 1234, AppTypes: []string{"expired"}}},
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := NewServer().QueryRegisteredServices(ctx, testCase.req)

			// Sort the received response according to the app instance id, to make the result
			// predictable
			sort.Slice(response.Infos, func(i, j int) bool {
				return response.Infos[i].AppInstanceId < response.Infos[j].AppInstanceId
			})

			if testCase.expected.err == nil {
				assert.Equal(t, testCase.expected.out, response.Infos)
				assert.NoError(t, testCase.expected.err)

			} else {
				assert.Equal(t, testCase.expected.err, err)
			}

			var sizeRegisteredServices int
			registeredServices.Range(func(k, v interface{}) bool {
				sizeRegisteredServices++
				return true
			})
			assert.Equal(t, sizeRegisteredServices, validService, "expiredServiceShouldBeDeleted")

		})
	}
}

func TestUnregisterService(t *testing.T) {

	eraseStorage()

	initialStorageList()

	type expectation struct {
		out *pb.UnregisterServiceResponse
		err error
	}
	testCases := []struct {
		name       string
		req        *pb.UnregisterServiceRequest
		entryExist bool
		expected   expectation
	}{
		{
			name:       "Smoketest",
			req:        &pb.UnregisterServiceRequest{Info: &pb.ServiceInfo{AppInstanceId: "driver-1"}},
			entryExist: true,
			expected: expectation{
				out: &pb.UnregisterServiceResponse{},
				err: nil,
			},
		},
		{
			name:       "AppInstanceIdNotFound",
			req:        &pb.UnregisterServiceRequest{Info: &pb.ServiceInfo{AppInstanceId: "not-found"}},
			entryExist: false,
			expected: expectation{
				out: &pb.UnregisterServiceResponse{},
				err: status.Error(codes.InvalidArgument, "AppInstanceId not found inside registered services"),
			},
		},
		{
			name:       "NoInfoValue",
			req:        &pb.UnregisterServiceRequest{Info: nil},
			entryExist: false,
			expected: expectation{
				out: &pb.UnregisterServiceResponse{},
				err: status.Error(codes.InvalidArgument, "No info{} value given. Nothing to delete."),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.entryExist {
				_, found := registeredServices.Load(testCase.req.Info.AppInstanceId)
				assert.True(t, found && testCase.entryExist)
			}

			ctx := context.Background()
			_, err := NewServer().UnregisterService(ctx, testCase.req)
			if testCase.expected.err == nil {
				assert.NoError(t, testCase.expected.err)
				_, found := registeredServices.Load(testCase.req.Info.AppInstanceId)
				assert.False(t, found)
			} else {
				assert.Equal(t, testCase.expected.err, err)
			}

		})
	}
}
