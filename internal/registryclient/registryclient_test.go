/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package registryclient

import (
	"context"
	"net"
	"testing"

	pb "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type mockRegistryServer struct {
	pb.UnimplementedRegistryApiServer
	registerCallCount   int
	unregisterCallCount int
	expireTime          uint32
	registerErr         error
	unregisterErr       error
}

func (m *mockRegistryServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	m.registerCallCount++
	if m.registerErr != nil {
		return nil, m.registerErr
	}
	return &pb.RegisterServiceResponse{ExpireTime: m.expireTime}, nil
}

func (m *mockRegistryServer) UnregisterService(ctx context.Context, req *pb.UnregisterServiceRequest) (*pb.UnregisterServiceResponse, error) {
	m.unregisterCallCount++
	if m.unregisterErr != nil {
		return nil, m.unregisterErr
	}
	return &pb.UnregisterServiceResponse{}, nil
}

func (m *mockRegistryServer) QueryRegisteredServices(ctx context.Context, req *pb.QueryRegisteredServicesRequest) (*pb.QueryRegisteredServicesResponse, error) {
	return &pb.QueryRegisteredServicesResponse{}, nil
}

func setupRegistryClientWithMock(t *testing.T, mock *mockRegistryServer) (*GrpcServerRegistry, func()) {
	lis := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	pb.RegisterRegistryApiServer(server, mock)

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	r := New("bufnet", "test-al", "localhost:9090")
	r.connection = conn
	r.client = pb.NewRegistryApiClient(conn)
	r.appInstanceId = "test-app-instance-id"

	cleanup := func() {
		if conn != nil {
			conn.Close()
		}
		server.Stop()
		lis.Close()
	}

	return r, cleanup
}

func TestConnect(t *testing.T) {
	t.Run("Connect with invalid address", func(t *testing.T) {
		r := New("invalid-address-without-port", "test-al-id", "localhost:9090")
		err := r.connect()

		assert.NoError(t, err, "connect() should not return an error")
		assert.NotNil(t, r.connection, "connection should not be nil")
		assert.NotNil(t, r.client, "client should not be nil")

		if r.connection != nil {
			r.connection.Close()
		}
	})

	t.Run("Connect with valid localhost address", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "localhost:9090")
		err := r.connect()

		assert.NoError(t, err, "connect() should not return an error")
		assert.NotNil(t, r.connection, "connection should not be nil")
		assert.NotNil(t, r.client, "client should not be nil")

		if r.connection != nil {
			r.connection.Close()
		}
	})

	t.Run("Connect with IPv4 address", func(t *testing.T) {
		r := New("127.0.0.1:8080", "al-001", "127.0.0.1:9090")
		err := r.connect()

		assert.NoError(t, err, "connect() should not return an error")
		assert.NotNil(t, r.connection, "connection should not be nil")
		assert.NotNil(t, r.client, "client should not be nil")

		if r.connection != nil {
			r.connection.Close()
		}
	})
}

func TestStop(t *testing.T) {
	t.Run("Stop without connection", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "localhost:9090")
		assert.Nil(t, r.connection, "New() should initialize with nil connection")
		r.Stop()
	})

	t.Run("Stop with active connection", func(t *testing.T) {
		mock := &mockRegistryServer{expireTime: 60}
		r, cleanup := setupRegistryClientWithMock(t, mock)
		defer cleanup()

		r.Stop()
		assert.Equal(t, 1, mock.unregisterCallCount, "unregister should be called once")
		assert.Nil(t, r.connection, "connection should be nil after Stop")
		assert.Nil(t, r.client, "client should be nil after Stop")
	})
}

func TestDisconnect(t *testing.T) {
	t.Run("Disconnect without connection initialized", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "localhost:9090")

		r.disconnect()
		assert.Nil(t, r.connection, "disconnect() should keep connection nil")
		assert.Nil(t, r.client, "disconnect() should keep client nil")
	})

	t.Run("Disconnect with active connection", func(t *testing.T) {
		mock := &mockRegistryServer{expireTime: 60}
		r, cleanup := setupRegistryClientWithMock(t, mock)
		defer cleanup()

		r.disconnect()
		assert.Equal(t, 1, mock.unregisterCallCount, "unregister should be called once")
		assert.Nil(t, r.connection, "connection should be nil after disconnect")
		assert.Nil(t, r.client, "client should be nil after disconnect")
	})
}

func TestRegister(t *testing.T) {
	t.Run("Register with valid parameters", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "localhost:9090")
		r.Register()
		assert.Equal(t, "test-al", r.alId, "alId should match")
	})

	t.Run("Register with IPv4 address", func(t *testing.T) {
		r := New("127.0.0.1:8080", "al-001", "127.0.0.1:9090")
		r.Register()
		assert.Equal(t, "al-001", r.alId, "alId should match")
	})

	t.Run("Register with DNS name", func(t *testing.T) {
		r := New("registry.example.com:8080", "al-dns", "service.example.com:9090")
		r.Register()
		assert.Equal(t, "al-dns", r.alId, "alId should match")
	})
}

func TestUnregister(t *testing.T) {
	t.Run("Unregister without connection panics", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "localhost:9090")
		assert.Panics(t, func() {
			_ = r.unregister()
		}, "unregister() should panic when client is nil")
	})

	t.Run("Unregister with connection succeeds", func(t *testing.T) {
		mock := &mockRegistryServer{expireTime: 60}
		r, cleanup := setupRegistryClientWithMock(t, mock)
		defer cleanup()

		err := r.unregister()
		assert.NoError(t, err, "unregister() should not return error with valid connection")
		assert.Equal(t, 1, mock.unregisterCallCount, "unregister should be called once")
	})

	t.Run("Unregister with server error returns error", func(t *testing.T) {
		mock := &mockRegistryServer{
			expireTime:    60,
			unregisterErr: assert.AnError,
		}
		r, cleanup := setupRegistryClientWithMock(t, mock)
		defer cleanup()

		err := r.unregister()
		assert.Error(t, err, "unregister() should return error when server fails")
		assert.Equal(t, 1, mock.unregisterCallCount, "unregister should be called once")
	})
}

func TestRegisterFunc(t *testing.T) {
	t.Run("register with invalid gRPC address", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "invalid-address")

		availableCSInterfaces = []string{INTERFACE_DRVINFO_V1}
		expireTime, err := r.register()
		assert.Error(t, err, "register() should return an error")
		assert.Equal(t, uint32(retryRegistrationInterval), expireTime, "register() expireTime should match expected")
	})

	t.Run("register with valid address but no server", func(t *testing.T) {
		r := New("localhost:8080", "test-al", "192.168.1.1:9090")

		availableCSInterfaces = []string{INTERFACE_DRVINFO_V1}
		expireTime, err := r.register()
		assert.Error(t, err, "register() should return an error")
		assert.Equal(t, uint32(retryRegistrationInterval), expireTime, "register() expireTime should match expected")
		expectedAppInstanceId := CDM_DEVICE_CLASS_DRIVER.String() + "-" + "test-al"
		assert.Equal(t, expectedAppInstanceId, r.appInstanceId, "appInstanceId should match expected")
		if r.connection != nil {
			r.connection.Close()
		}
	})

	t.Run("register with mock server succeeds", func(t *testing.T) {
		mock := &mockRegistryServer{expireTime: 120}
		r, cleanup := setupRegistryClientWithMock(t, mock)
		defer cleanup()

		availableCSInterfaces = []string{INTERFACE_DRVINFO_V1}
		r.appInstanceId = CDM_DEVICE_CLASS_DRIVER.String() + "-" + "test-al"

		register := &pb.RegisterServiceRequest{Info: &pb.ServiceInfo{
			AppTypes:         getCsInterfaces(),
			Interfaces:       getCsInterfaces(),
			AppInstanceId:    r.appInstanceId,
			DriverSchemaUris: []string{r.alId},
			GrpcIp:           &pb.ServiceInfo_Ipv4Address{Ipv4Address: "192.168.1.1"},
			GrpcIpPortNumber: 9090,
		}}

		response, err := r.client.RegisterService(context.Background(), register)

		assert.NoError(t, err, "RegisterService should not return error with mock server")
		assert.Equal(t, uint32(120), response.ExpireTime, "register() should return correct expireTime")
		assert.Equal(t, 1, mock.registerCallCount, "RegisterService should be called once")
	})
}
