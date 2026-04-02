---
title: "Troubleshooting"
nav_order: 6
---

## Troubleshooting

This guide provides solutions for common issues encountered when working with Asset Links,
particularly focusing on gRPC connection failures and discovery agent problems.

## gRPC Connection Failures

### Overview

Asset Links communicate with the Asset Gateway through gRPC connections.
An Asset Link acts  as a **gRPC server** for all its exposed services (e.g., discovery) and as a **gRPC client** towards the registry.
Two primary connection types exist:

- **Client Connection** (e.g., Registry): Asset Link → gRPC Server Registry (example: port 50051)
- **Server Connection** (e.g., Discovery): Asset Gateway/al-ctl → Asset Link (example: port 8081)

> **Note**: Port numbers may vary as per the configuration. Examples use common default ports but can be adjusted accordingly.

### Common Error Messages

#### **Error**: "Could not dial grpc registry"

**Cause**: The Asset Link cannot establish a connection to the gRPC Server Registry.

**Possible Reasons**:

- Registry server is not running
- Incorrect registry address or port
- Network connectivity issues
- Firewall blocking the connection

**Solutions**:

1. **Verify Registry Server is Running**:

   ```bash
   # Check if registry container is running
   docker ps | grep grpc-server-registry
   
   # Start the registry server
   docker-compose -f registry/docker-compose.yml up
   ```

2. **Check Registry Address Configuration**:

   Verify that the `--grpc-server-address` command-line argument is set correctly when starting the Asset Link:

   ```bash
   # When running on Windows using docker desktop
   --grpc-server-address=localhost:<PORT>
   
   # When running on Linux
   --grpc-server-address=$(hostname):<PORT>
   ```

   Additionally, check the configured address in `docker-compose.yml` environment variables, or use `docker ps` to see which ports are exposed by running containers.

3. **Test Network Connectivity**:

   ```bash
   # Test if port is accessible (Linux/macOS)
   # Replace <PORT> with configured registry port (e.g., 50051)
   nc -zv localhost <PORT>
   
   # Windows PowerShell
   Test-NetConnection -ComputerName localhost -Port <PORT>
   ```

4. **Check Firewall Settings**: Ensure configured port is not blocked by firewall.

#### **Error**: "can not connect with server"

**Cause**: The discovery client (al-ctl or Asset Gateway) cannot connect to the Asset Link server.

**Possible Reasons**:

- Asset Link server is not running
- Incorrect endpoint address
- Port already in use
- Server crashed or failed to start

**Solutions**:

1. **Verify Asset Link is Running**:

   Check the Asset Link logs for startup messages. A successful start will show the gRPC listen address:

   ```bash
   # Look for a log like:
   # "Serving gRPC server <address>"
   # Example: Serving gRPC Server address=localhost:8080
   ```

   If this message is not visible, the server failed to start. Review earlier log entries for errors.

2. **Check Endpoint Configuration**:

   ```bash
   # When testing with al-ctl (replace <PORT> with Asset Link port, e.g., 8081)
   al-ctl discover -e localhost:<PORT>
   
   # Verify the Asset Link is listening on the correct port
   netstat -an | grep <PORT>  # Linux/macOS
   netstat -an | findstr <PORT>  # Windows
   ```

   In the `netstat` output, look for a line showing `LISTEN` state on the expected port. If the port does not appear, the Asset Link is not running or is bound to a different port.

3. **Verify Port Availability**:

   ```bash
   # Check if another process is using the port (replace <PORT> with configured port)
   lsof -i :<PORT>  # Linux/macOS
   netstat -ano | findstr :<PORT>  # Windows
   ```

   If the output shows a process (PID) already bound to the port, another application is occupying it. Either stop that process or configure the Asset Link to use a different port.

#### **Error**: "Could not parse GRPC server address"

**Cause**: The gRPC server address format is invalid.

**Expected Format**: `host:port` (e.g., `localhost:8081`, `192.168.1.100:8081`)

**Invalid Examples**:

- `8081` (missing host)
- `localhost` (missing port)
- `localhost:abc` (non-numeric port)
- `http://localhost:8081` (includes protocol)

> **Note**: Replace port numbers with actual configured values.

**Solution**: Use the correct format without protocol prefix.

---

## Discovery Agent Issues

### "No Discovery implementation found"

**Error Code**: `UNIMPLEMENTED (12)`

**Cause**: The Asset Link does not implement the Discovery interface.

**Solution**: Ensure Asset Link implements the required Discovery interface methods:

- `Discover(discoveryConfig config.DiscoveryConfig, output chan<- *generated.DiscoverResponse) error`
- `GetSupportedFilters() ([]*generated.FilterType, error)`
- `GetSupportedOptions() ([]*generated.OptionType, error)`

For a working example, see the reference implementation in [`cdm-al-reference/reference/reference.go`](https://github.com/industrial-asset-hub/asset-link-sdk/tree/main/cdm-al-reference/reference/reference.go).

### Discovery Timeout Issues

**Symptoms**:

- Discovery operations hang indefinitely
- No devices found when devices should be present
- Context cancellation errors

**Solutions**:

1. **Set Appropriate Timeout**:

   ```bash
   # Set a 30-second timeout (replace <PORT> with Asset Link port)
   al-ctl discover -e localhost:<PORT> -n 30
   
   # Use fractional seconds for quick tests
   al-ctl discover -e localhost:<PORT> -n 0.5
   ```

2. **Increase Timeout for Large Networks**:

   - For networks with many devices, use longer timeouts (60-300 seconds)
   - Discovery time scales with number of devices and network latency

3. **Check Discovery Configuration**:

   - Review filters and options in discovery configuration file
   - Overly broad filters may cause excessive scanning

### Discovery Returns No Devices

**Possible Causes**:

1. **Incorrect Filters**: Filters are too restrictive
2. **Network Issues**: Devices are unreachable
3. **Authentication Failures**: Invalid credentials
4. **Protocol Issues**: Device communication protocol mismatch

**Troubleshooting Steps**:

1. **Test with Minimal Filters**:

   ```json
   {
     "filters": [],
     "options": []
   }
   ```

2. **Check Discovery Logs**: Look for specific error codes in the response:

   ```bash
   # Replace <PORT> with Asset Link port
   al-ctl discover -e localhost:<PORT> --log-level debug
   ```

3. **Other Issues**: For network, authentication, or protocol problems, test device connectivity independently, verify credentials, confirm the correct protocol is in use, and capture network traffic with a tool like Wireshark to identify the root cause.

---

## Error Codes Reference

### Standard gRPC Error Codes

| Code | Value | Description | Common Causes |
| --- | --- | --- | --- |
| `OK` | 0 | Success | Operation completed successfully |
| `CANCELLED` | 1 | Operation cancelled | User cancellation or timeout |
| `UNKNOWN` | 2 | Unknown error | Unexpected internal error |
| `INVALID_ARGUMENT` | 3 | Invalid request | Malformed discovery config, invalid parameters |
| `DEADLINE_EXCEEDED` | 4 | Timeout expired | Discovery took too long, increase timeout |
| `NOT_FOUND` | 5 | Not found | Device or resource doesn't exist |
| `ALREADY_EXISTS` | 6 | Already exists | Duplicate registration |
| `PERMISSION_DENIED` | 7 | Access denied | Insufficient permissions |
| `RESOURCE_EXHAUSTED` | 8 | Resource limit | Discovery already running, system overloaded |
| `FAILED_PRECONDITION` | 9 | Invalid state | Missing configuration, not ready |
| `ABORTED` | 10 | Aborted | Operation conflict |
| `OUT_OF_RANGE` | 11 | Out of range | Parameter value invalid |
| `UNIMPLEMENTED` | 12 | Not implemented | Feature not available |
| `INTERNAL` | 13 | Internal error | Driver failure, unexpected exception |
| `UNAVAILABLE` | 14 | Service unavailable | Server down, network issue (retriable) |
| `DATA_LOSS` | 15 | Data corruption | Unrecoverable data error |
| `UNAUTHENTICATED` | 16 | Not authenticated | Missing or invalid credentials |

### Custom Error Codes

| Code | Value | Description |
| --- | --- | --- |
| `INVALID_DATAPOINT_ID` | 64 | Invalid datapoint identifier |
| `INVALID_SUBSCRIPTION_ID` | 65 | Invalid subscription identifier |
| `PAYLOAD_SERIALISATION_ISSUE` | 66 | Failed to serialize payload |
| `PAYLOAD_DESERIALISATION_ISSUE` | 67 | Failed to deserialize payload |
| `BUSY` | 68 | System busy, operation in progress |

---

## Configuration Issues

### Port Conflicts

**Common Default Ports** (configurable):

- gRPC Server: `8081` (example)
- HTTP Server (Observability): `8082` (example)
- gRPC Registry: `50051` (example)
- Visualization Server: `8083` (example)

> **Note**: These are commonly used default ports. Asset Link may use different ports configured.

**Symptoms**: Asset Link fails to start with "address already in use" error.

**Solutions**:

1. **Change Port Configuration**:

   ```bash
   # Use different gRPC port (replace <NEW_PORT> with desired port number)
   ./asset-link --grpc-server-address=:<NEW_PORT>
   ```

2. **Kill Conflicting Process**:

   ```bash
   # Find process using the port (replace <PORT> with port number)
   lsof -i :<PORT>  # Linux/macOS
   netstat -ano | findstr :<PORT>  # Windows
   
   # Kill the process (use PID from above)
   kill -9 <PID>  # Linux/macOS
   taskkill /PID <PID> /F  # Windows
   ```

### Registration Issues

**Symptom**: Asset Link starts but does not appear in gateway.

**Possible Causes**:

1. Registry address misconfigured
2. Registration refresh failing
3. Network issues between Asset Link and registry

**Log Messages to Check**:

```text
INFO: Register asset link at grpc server registry <grpc-registry-address>

For successful Registration, check below logs:
INFO:Serving gPRC Server <grpc-server-address>
INFO:GetVersionInfo called

If unsuccessful, below logs will be visible
WARN: Could not register at grpc server registry
```

`INFO` indicates that registration was attempted.
`WARN` message indicates a failure to check the accompanying error details for the root cause (e.g., wrong address, network unreachable).

**Solutions**:

1. **Check Registration Address Parameter**: Ensure `--grpc-server-endpoint-address` is correct

   ```bash
   # For Windows Docker Desktop only, use the hostname:
   --grpc-server-endpoint-address=host.docker.internal
   
   # For Windows (non-Docker), use the machine's IP address. For example:
   --grpc-server-endpoint-address=192.168.1.100

   # For Linux, use hostname as below:
   --grpc-server-endpoint-address=$(hostname)
   ```

2. **Review Network Configuration**: Ensure registry can reach Asset Link's advertised address

3. **Check Registration Intervals**:

   - Retry interval on error: 10 seconds
   - Re-registration refresh: 60 seconds

---

## Docker and Docker Compose Issues

### Container Networking

**Issue**: Services can't communicate despite being in same compose file.

**Solution**: Verify network configuration in docker-compose.yml:

```yaml
services:
  asset-link:
    networks:
      - cdm
  grpc-server-registry:
    networks:
      - cdm

networks:
  cdm:
    name: cdm
```

### Container Not Starting

**Steps**:

1. Check container logs for error messages, stack traces, or failed health checks:

   ```bash
   docker logs <container-name>
   ```

   Look for `ERROR` or `FATAL` level messages that indicate why the container failed to start.

2. Verify image availability:

   ```bash
   docker images | grep <container-name>
   # e.g.
   docker images | grep grpc-server-registry
   ```

3. Check resource constraints (memory, CPU)

---

## Debugging Tips

### Enable Debug Logging

```bash
# Asset Link server
./asset-link --log-level=debug

# al-ctl tool (replace <PORT> with Asset Link port)
al-ctl discover -e localhost:<PORT> --log-level=debug
```

### Verify gRPC Communication

Use `grpcurl` to test endpoints directly:

```bash
# Replace <GRPC_PORT> with Asset Link gRPC port (e.g., 8081)
# List available services
grpcurl -plaintext localhost:<GRPC_PORT> list

# Describe a service
grpcurl -plaintext localhost:<GRPC_PORT> describe siemens.industrialassethub.discover.v1.DeviceDiscoverApi

# Test discovery
grpcurl -plaintext -d '{"filters":[],"options":[]}' \
  localhost:<GRPC_PORT> siemens.industrialassethub.discover.v1.DeviceDiscoverApi/DiscoverDevices
```

### Monitor Network Traffic

```bash
# Linux/macOS - capture gRPC traffic (replace with configured ports)
tcpdump -i any -A 'tcp port <GRPC_PORT> or tcp port <REGISTRY_PORT>'

# Windows - use Wireshark with filter (replace with configured ports)
tcp.port == <GRPC_PORT> || tcp.port == <REGISTRY_PORT>
```

### Check Service Health

```bash
# HTTP observability endpoint (replace <HTTP_PORT> with configured observability port, e.g., 8082)
curl http://localhost:<HTTP_PORT>/health
curl http://localhost:<HTTP_PORT>/version
curl http://localhost:<HTTP_PORT>/stats
```

---

## Reference Documentation

For additional information, refer to:

- **Asset Gateway**: [https://github.com/industrial-asset-hub/asset-gateway](https://github.com/industrial-asset-hub/asset-gateway)
  - Contains the gRPC Server Registry implementation
  - Gateway-side discovery client implementation
  - Complete system architecture documentation

- **Asset Link SDK Overview**: [Overview](overview.md)
- **Bootstrapping Guide**: [Bootstrap](bootstrap.md)
- **Command Line Tool**: [Command Line Tool](cmd.md)
- **Model Documentation**: [Working with the Asset Model](model.md)

### Protocol Buffer Specifications

The gRPC interfaces are defined in the `specs/` directory:

- `conn_suite_registry.proto` - Registry service interface
- `conn_suite_drv_info.proto` - Driver information interface
- `iah_discover.proto` - Discovery service interface
- `common_identifiers.proto` - Identifiers service interface

---

## Getting Help

If issues persist:

1. **Check Logs**: Enable debug logging and review output carefully
2. **Review Configuration**: Verify all addresses, ports, and parameters
3. **Test Connectivity**: Use network tools to verify basic connectivity
4. **Consult Error Codes**: Match error codes to this guide for specific solutions
5. **Open an Issue**: [GitHub Issues](https://github.com/industrial-asset-hub/asset-link-sdk/issues)

When reporting issues, include:

- Asset Link SDK version
- Asset Gateway version (if applicable)
- Full error messages and logs
- Configuration files (redact sensitive information)
- Network topology (Docker, bare metal, cloud)
- Steps to reproduce
