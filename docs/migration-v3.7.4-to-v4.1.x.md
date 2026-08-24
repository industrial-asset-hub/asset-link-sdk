---
title: "Migration Guide: v3.7.4 to v4.1.x"
nav_order: 7
---

# Migration Guide: v3.7.4 → v4.1.x

This guide covers all breaking changes and required code updates when migrating an Asset Link implementation from SDK **v3.7.4** (module `github.com/industrial-asset-hub/asset-link-sdk/v3`) to **v4.1.x** (module `github.com/industrial-asset-hub/asset-link-sdk/v4`).

---

## Table of Contents

1. [Module path change (v3 → v4)](#1-module-path-change-v3--v4)
2. [Identifiers API replaced by DeviceInfo API](#2-identifiers-api-replaced-by-deviceinfo-api)
3. [Asset data model changes (base schema v0.12.0 → v1.19.0)](#3-asset-data-model-changes-base-schema-v0120--v1180)
4. [Builder API changes](#4-builder-api-changes)
5. [Metadata struct: new optional fields](#5-metadata-struct-new-optional-fields)
6. [Model helper method signature changes (errors are now returned)](#6-model-helper-method-signature-changes-errors-are-now-returned)
7. [Removed types and functions](#7-removed-types-and-functions)
8. [New types and functions](#8-new-types-and-functions)
9. [Registry interface constant rename](#9-registry-interface-constant-rename)
10. [New generated package: conn_suite_device_info](#10-new-generated-package-conn_suite_device_info)
11. [Dependency updates](#11-dependency-updates)

---

## 1. Module path change (v3 → v4)

All import paths must be updated from `v3` to `v4`.

**Before:**
```go
import "github.com/industrial-asset-hub/asset-link-sdk/v3/assetlink"
import "github.com/industrial-asset-hub/asset-link-sdk/v3/model"
import "github.com/industrial-asset-hub/asset-link-sdk/v3/config"
import "github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"
import generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
```

**After:**
```go
import "github.com/industrial-asset-hub/asset-link-sdk/v4/assetlink"
import "github.com/industrial-asset-hub/asset-link-sdk/v4/model"
import "github.com/industrial-asset-hub/asset-link-sdk/v4/config"
import "github.com/industrial-asset-hub/asset-link-sdk/v4/metadata"
import generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
```

Update `go.mod`:
```
module your-asset-link

go 1.25.13

require (
    github.com/industrial-asset-hub/asset-link-sdk/v4 <latest-version>
)
```

Run `go mod tidy` afterwards.

---

## 2. Identifiers API replaced by DeviceInfo API

The `Identifiers` interface and the gRPC `IdentifiersApi` have been **removed** and replaced by the new **DeviceInfo** API (`conn_suite_device_info`).

### Interface change

**Before** (`internal/features/features.go`):
```go
// config.IdentifiersRequest is a helper wrapping GetIdentifiersRequest
type Identifiers interface {
    GetIdentifiers(identifiersRequest config.IdentifiersRequest) ([]*generated.DeviceIdentifier, error)
}
```

**After:**
```go
import deviceinfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"

type DeviceInfo interface {
    GetPropertyValues(request *deviceinfo.GetPropertyValuesRequest) (*deviceinfo.GetPropertyValuesResponse, error)
    GetSupportedProperties(request *deviceinfo.GetSupportedPropertiesRequest) (*deviceinfo.GetSupportedPropertiesResponse, error)
}
```

### Implementation change

**Before:**
```go
func (m *MyAssetLink) GetIdentifiers(req config.IdentifiersRequest) ([]*generated.DeviceIdentifier, error) {
    // build and return identifiers
    return identifiers, nil
}
```

**After:**

The `GetPropertyValues` method receives a device target in the request and must return the asset's properties as a flat list of key/value pairs. Use `DeviceInfo.ConvertToPropertyValueResults()` to convert a populated `DeviceInfo` into the expected response format:

```go
import (
    deviceinfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
    generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
)

func (m *MyAssetLink) GetPropertyValues(req *deviceinfo.GetPropertyValuesRequest) (*deviceinfo.GetPropertyValuesResponse, error) {
    // 1. Extract connection parameters from the request target
    device := req.GetDevice()
    if device == nil {
        return nil, status.Errorf(codes.InvalidArgument, "missing device target")
    }
    paramJSON := device.GetConnectionParameterSet().GetParameterJson()
    credentials := device.GetConnectionParameterSet().GetCredentials()

    // 2. Retrieve device details using paramJSON / credentials (protocol-specific)
    deviceDetails, err := m.retrieveDevice(paramJSON, credentials)
    if err != nil {
        return nil, status.Errorf(codes.Unavailable, "could not reach device: %v", err)
    }

    // 3. Build the DeviceInfo model
    deviceInfo, err := buildDeviceInfo(deviceDetails)
    if err != nil {
        return nil, err
    }

    // 4. Convert to property value results using the SDK helper
    results, err := deviceInfo.ConvertToPropertyValueResults()
    if err != nil {
        return nil, err
    }
    return &deviceinfo.GetPropertyValuesResponse{PropertyResults: results}, nil
}

func (m *MyAssetLink) GetSupportedProperties(_ *deviceinfo.GetSupportedPropertiesRequest) (*deviceinfo.GetSupportedPropertiesResponse, error) {
    return &deviceinfo.GetSupportedPropertiesResponse{
        Properties: []*deviceinfo.SupportedProperty{
            {Key: "name", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
            {Key: "functional_object_type", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
            {Key: "functional_object_schema_url", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
            {Key: "asset_identifiers", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
            {Key: "connection_points", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
            {Key: "software_components", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
            {Key: "product_instance_information", Type: &deviceinfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRUCT}},
        },
    }, nil
}
```

`ConvertToPropertyValueResults()` internally calls `deviceInfo.ConvertToJson()` (which produces a `map[string]interface{}`) and then wraps each key/value pair into a `PropertyValueResult`. You do not need to call `ConvertToJson` manually.

See `cdm-al-reference/reference/reference.go` for a complete working example.

### Removed config helper

`config.IdentifiersRequest`, `config.NewIdentifiersRequestFromGetIdentifiersReq`, and `config/identifiersrequest_impl.go` have been **deleted**. The new `config.DeviceInfoRequest` interface (`config/deviceinforequest.go`) exists but is not directly used by the `DeviceInfo` interface — requests arrive as typed protobuf messages from the `conn_suite_device_info` package.

### registry.json change

Replace `"siemens.common.identifiers.v1"` with `"siemens.connectivitysuite.deviceinfo.v1"` in your `registry.json`. Keep the other interface entries intact:

**Before:**
```json
{
  "app_instance_id": "cdm-device-class-driver-com.example.al.mylink",
  "app_types": [
    "siemens.connectivitysuite.drvinfo.v1",
    "siemens.industrialassethub.discover.v1",
    "siemens.common.identifiers.v1"
  ],
  "driver_schema_uris": ["com.example.al.mylink"]
}
```

**After:**
```json
{
  "app_instance_id": "cdm-device-class-driver-com.example.al.mylink",
  "app_types": [
    "siemens.connectivitysuite.drvinfo.v1",
    "siemens.industrialassethub.discover.v1",
    "siemens.connectivitysuite.deviceinfo.v1"
  ],
  "driver_schema_uris": ["com.example.al.mylink"]
}
```

If your Asset Link does not implement the DeviceInfo interface, simply omit `"siemens.connectivitysuite.deviceinfo.v1"` from `app_types`.

---

## 3. Asset data model changes (base schema v0.12.0 → v1.19.0)

The base schema has been updated from `v0.12.0` to `v1.19.0`. This is a significant restructuring of the asset model.

### 3.1 Schema URL

| Before | After |
|---|---|
| `https://schema.industrial-assets.io/base/v0.12.0` | `https://industrial-assets.io/schemas/iah/base-schema/released/v1` |

The canonical schema URL is now exposed as the constant `model.FunctionalObjectSchemaUrl`:
```go
const FunctionalObjectSchemaUrl = "https://industrial-assets.io/schemas/iah/base-schema/released/v1/iah-base.json"
```

### 3.2 `DeviceInfo` struct

| Removed field | Replacement |
|---|---|
| `Type string` (`@type`) | `FunctionalObjectType any` (`functional_object_type`) |
| `Context *AssetContext` (`@context`) | `FunctionalObjectSchemaUrl any` (`functional_object_schema_url`) |
| `MacIdentifiers []MacIdentifier` | `AssetIdentifiers []any` (use `AddMacIdentifier`) |

### 3.3 `Asset` struct (embedded in `DeviceInfo`)

| Removed field | Replacement / note |
|---|---|
| `Id string` | Removed (no longer required) |
| `ManagementState ManagementState` | Removed — no longer part of the schema |
| `ReachabilityState *ReachabilityState` | Removed |
| `ProductInstanceIdentifier *ProductSerialIdentifier` | `ProductInstanceInformation interface{}` |
| `CustomUiProperties []CustomProperty` | Removed |
| `FunctionalParts []interface{}` | Removed |
| `LastModifiedTimestamp *time.Time` | Removed |
| `Zone *string` | Removed |
| `Responsible *string` | Removed |
| `OtherStates []State` | Removed |
| `AssetIdentifiers []interface{}` (omitempty) | `AssetIdentifiers []interface{}` (now **required**, not omitempty) |
| — | `AssetRelations []AssetRelation` (new) |
| — | `OperatingMode interface{}` (new) |

### 3.4 `GatewayInfo` / `NewGateway`

`GatewayInfo` and `NewGateway()` have been **removed**. Use the new `FunctionalObjectType` constants directly when needed:
```go
const GatewayFunctionalObjectTypeGateway GatewayFunctionalObjectType = "Gateway"
```

### 3.5 `NewDevice` now returns an error

**Before:**
```go
deviceInfo := model.NewDevice("EthernetDevice", assetName)
```

**After:**
```go
deviceInfo, err := model.NewDevice(string(model.DeviceFunctionalObjectTypeDevice), assetName)
if err != nil {
    // handle validation error
}
```

`NewDevice` now validates that the `functionalObjectType` is non-empty and belongs to the allowed set:
- `model.AssetFunctionalObjectTypeAsset` (`"Asset"`)
- `model.DeviceFunctionalObjectTypeDevice` (`"Device"`)
- `model.GatewayFunctionalObjectTypeGateway` (`"Gateway"`)
- `model.SoftwareArtifactFunctionalObjectTypeSoftwareArtifact` (`"SoftwareArtifact"`)

---

## 4. Builder API changes

### `Identifiers(...)` renamed to `DeviceInfo(...)`

**Before:**
```go
alImpl := assetlink.New(metadata.Metadata{...}).
    Discovery(myImpl).
    Identifiers(myImpl).
    Build()
```

**After:**
```go
alImpl := assetlink.New(metadata.Metadata{...}).
    Discovery(myImpl).
    DeviceInfo(myImpl).   // renamed
    Build()
```

If your Asset Link does not implement the DeviceInfo interface, simply omit the `.DeviceInfo(...)` call.

---

## 5. Metadata struct: new optional fields

Three new optional fields were added to `metadata.Metadata`:

**Before:**
```go
metadata.Metadata{
    AlId:    "com.example.al.mylink",
    AlName:  "My Asset Link",
    Vendor:  "My Company",
    Version: metadata.Version{Version: version, Commit: commit, Date: date},
}
```

**After:**
```go
metadata.Metadata{
    AlId:        "com.example.al.mylink",
    AlName:      "My Asset Link",
    Vendor:      "My Company",
    Version:     metadata.Version{Version: version, Commit: commit, Date: date},
    Description: "Short description of what this Asset Link does.",   // new
    DocUrl:      "https://example.com/docs",                          // new
    FeedbackUrl: "https://example.com/feedback",                      // new
}
```

All three new fields are optional strings and can be left empty.

---

## 6. Model helper method signature changes (errors are now returned)

All model helper methods that previously silently ignored invalid inputs now return an `error`. Callers **must** handle the returned errors.

### `AddNameplate`

**Before:**
```go
// parameter: manufacturerProductDesignation
deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productName, hardwareVersion, serialNumber)
```

**After:**
```go
// parameter renamed: productFamily (was manufacturerProductDesignation)
err := deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productFamily, hardwareVersion, serialNumber)
if err != nil {
    // handle error
}
```

Note the parameter rename: `manufacturerProductDesignation` → `productFamily`.

### `AddDescription`

**Before:**
```go
deviceInfo.AddDescription("My Device")
```

**After:**
```go
err := deviceInfo.AddDescription("My Device")
```

### `AddCapabilities`

**Before:**
```go
deviceInfo.AddCapabilities("firmware_update", false)
```

**After:**
```go
err := deviceInfo.AddCapabilities("firmware_update", false)
```

### `AddNic`

**Before:**
```go
nicId := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03")
```

**After:**
```go
nicId, err := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03")
if err != nil {
    // handle validation error (e.g., invalid MAC format)
}
```

`AddNic` now validates the MAC address against `MacAddressPattern` (`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`).

A new variant `AddNicWithoutMacIdentifier` is available when you want to add a NIC without automatically appending a `MacIdentifier` to `AssetIdentifiers`.

### `AddIPv4`

**Before:**
```go
id := deviceInfo.AddIPv4(nicId, "192.168.0.10", "255.255.255.0", "")
```

**After:**
```go
// err is nil here because IP and mask are non-empty (error fires only when ALL fields are empty)
id, err := deviceInfo.AddIPv4(nicId, "192.168.0.10", "255.255.255.0", "")
if err != nil {
    // only reached when ipv4Address, networkMask, AND routerAddress are all empty
}
```

**Important behavior change:** The error is returned only when **all** three address arguments are empty strings. If at least one field is non-empty, the call succeeds; invalid individual fields log a warning and are skipped. In v3 an empty `routerAddress` was silently ignored; in v4 the same call behaves identically.

### `AddIPv6`

**Before:**
```go
id := deviceInfo.AddIPv6(nicId, ipv6Addr, prefix, router)
```

**After:**
```go
id, err := deviceInfo.AddIPv6(nicId, ipv6Addr, prefix, router)
```

The same OR-logic applies: error only when all three fields are empty.

### `AddSoftware` renamed to `AddSoftwareArtifactComponent`

**Before:**
```go
deviceInfo.AddSoftware("Firmware", "1.0.0", true)
```

**After:**
```go
err := deviceInfo.AddSoftwareArtifactComponent("Firmware", "1.0.0", true)
```

A new `AddRunningSoftwareComponent` method is also available for running software instances identified by a `runningSoftwareId`.

---

## 7. Removed types and functions

| Removed symbol | Notes |
|---|---|
| `model.GatewayInfo` | Use `DeviceInfo` with `GatewayFunctionalObjectTypeGateway` |
| `model.NewGateway()` | Removed with `GatewayInfo` |
| `model.AssetContext` | JSON-LD `@context` is no longer part of the schema |
| `model.ArtifactChecksum` | Removed from base schema |
| `model.AssetIdentifier` / `AssetIdentifierAssetIdentifierType` | Removed; use specific identifier types |
| `model.ProductSerialIdentifier` | Replaced by `ProductInstanceInformation` |
| `model.ManagementState` / `ManagementStateValues*` | Removed from schema |
| `model.ReachabilityState` / `ReachabilityStateValues*` | Removed from schema |
| `model.State` | Removed from schema |
| `model.CustomProperty` | Removed from schema |
| `(*DeviceInfo).AddManagementState()` | Removed |
| `(*DeviceInfo).addReachabilityState()` | Removed (was unexported) |
| `(*DeviceInfo).addIdentifier()` | Replaced by `AddMacIdentifier()` (now exported) |
| `(*DeviceInfo).AddSoftware()` | Renamed to `AddSoftwareArtifactComponent()` |
| `model.getAssetContext()` | Removed |
| `model.getAssetCreationTimestamp()` | Removed |
| `model.GatewayHelper*` (`gateway-helper.go`) | File deleted |
| `config.IdentifiersRequest` | Replaced by `config.DeviceInfoRequest` |
| `config.NewIdentifiersRequestFromGetIdentifiersReq()` | Removed |
| `features.Identifiers` interface | Replaced by `features.DeviceInfo` |
| `internal/server/devicediscovery/identifiers.go` | Replaced by `deviceinfo.go` |
| `generated/conn_suite_drv_info` (IdentifiersApiServer) | No longer embedded in `alFeatureBuilder` |

---

## 8. New types and functions

| New symbol | Description |
|---|---|
| `model.ValidationError` | Returned when a field value fails pattern validation |
| `model.EmptyError` | Returned when a required field is empty |
| `model.PermissibleValuesError` | Returned when a value is not in the allowed set |
| `model.ErrValidation`, `model.ErrEmpty` | Sentinel values for use with `errors.Is` (see below) |
| `model.FunctionalObjectSchemaUrl` | Constant: canonical schema URL for IAH base schema v1 |
| `model.MacAddressPattern` | Regex pattern for MAC address validation |
| `model.IPv4AddressPattern` | Regex pattern for IPv4 address validation |
| `model.NetworkMaskPattern` | Regex pattern for network mask validation |
| `model.IPv6AddressPattern` | Regex pattern for IPv6 address validation |
| `model.IPv6NetworkPrefixPattern` | Regex pattern for IPv6 network prefix (CIDR notation) |
| `model.RouterIPv4AddressPattern` | Regex pattern for IPv4 router/gateway addresses |
| `model.RouterIPv6AddressPattern` | Regex pattern for IPv6 router/gateway addresses |
| `model.IdLinkPattern` | Regex pattern for IEC 61406 product links |
| `model.CustomIdentifierValuePattern` | Regex pattern for custom identifier values |
| `model.PredicatePattern` | Regex pattern for asset relation predicates |
| `(*DeviceInfo).AddMacIdentifier(mac)` | Exported helper (previously `addIdentifier`, unexported) |
| `(*DeviceInfo).AddIdLinkIdentifier(uri)` | Add an IdLink identifier |
| `(*DeviceInfo).AddCustomIdentifier(name, value)` | Add a named custom identifier |
| `(*DeviceInfo).AddCertificateIdentifier(certID)` | Add a certificate identifier |
| `(*DeviceInfo).AddHostBasedSoftwareIdentifier(name, version, hostIdentifier)` | Add a host-based software identifier |
| `(*DeviceInfo).AddParentRelativeIdentifier(parentIdentifier, slot, subslot)` | Add a slot/subslot identifier relative to a parent asset |
| `(*DeviceInfo).AddAssetRelation(...)` | Add a relationship to another asset |
| `(*DeviceInfo).AddNicWithoutMacIdentifier(name, mac)` | Add a NIC without auto-adding a MAC identifier |
| `(*DeviceInfo).AddRunningSoftwareComponent(name, version, isFirmware, runningSoftwareId)` | Add a running software component (instance identified by `runningSoftwareId`) |
| `(*DeviceInfo).AddProductInstanceIdentifier(vendor, articleNumber, serialNumber)` | Add a vendor/article-number/serial-number identifier; skipped only when all three fields are empty |
| `model.AssetRelation` | New type representing a directed asset relationship |
| `model.RelatedAsset` | Asset stub used in `AssetRelation` |
| `model.ProductInstanceInformation` | Replaces `ProductSerialIdentifier` |
| `features.DeviceInfo` interface | Replaces `features.Identifiers` |
| `generated/conn_suite_device_info` package | New gRPC package for the DeviceInfo API |
| `(*DeviceInfo).ConvertToPropertyValueResults()` | Converts a `DeviceInfo` to `[]*PropertyValueResult` for `GetPropertyValues` responses |
| `(*DeviceInfo).ConvertToJson()` | Marshals `DeviceInfo` to `map[string]interface{}`; used internally by `ConvertToPropertyValueResults` |

### Error type usage with `errors.Is` and `errors.As`

Use `errors.Is` with the provided sentinels to check the error category without needing the details:

```go
import (
    "errors"
    "github.com/industrial-asset-hub/asset-link-sdk/v4/model"
)

deviceInfo, err := model.NewDevice(string(model.DeviceFunctionalObjectTypeDevice), name)
if err != nil {
    switch {
    case errors.Is(err, model.ErrEmpty):
        log.Warn().Msgf("required field empty: %v", err)
    case errors.Is(err, model.ErrValidation):
        log.Warn().Msgf("validation failed: %v", err)
    default:
        return err
    }
}
```

Use `errors.As` when you need the structured error fields:

```go
deviceInfo, err := model.NewDevice(string(model.DeviceFunctionalObjectTypeDevice), name)
if err != nil {
    var emptyErr *model.EmptyError
    var permErr *model.PermissibleValuesError
    switch {
    case errors.As(err, &emptyErr):
        log.Warn().Msgf("required field empty: %s", emptyErr.Field)
    case errors.As(err, &permErr):
        log.Warn().Msgf("invalid value for %s: %v", permErr.Field, permErr.Value)
    default:
        return err
    }
}
```

Note: `model.ErrPermissibleValues` is not defined; use `errors.As` with `*model.PermissibleValuesError` to match that type.

---

## 9. Registry interface constant rename

**Before:**
```go
registryclient.INTERFACE_COMMON_IDENTIFIERS_V1 // "siemens.common.identifiers.v1"
```

**After:**
```go
registryclient.INTERFACE_CONN_SUITE_DEVICEINFO_V1 // "siemens.connectivitysuite.deviceinfo.v1"
```

---

## 10. New generated package: conn_suite_device_info

A new gRPC package has been added at `generated/conn_suite_device_info`. It contains the protobuf-generated code for the new DeviceInfo API including:

- `DeviceInfoApiServer` / `DeviceInfoApiClient` interfaces
- `GetPropertyValuesRequest` / `GetPropertyValuesResponse`
- `GetSupportedPropertiesRequest` / `GetSupportedPropertiesResponse`
- `SupportedProperty` with typed `Datatype` field

The proto definitions are in `specs/conn_suite_device_info.proto` and `specs/common_properties.proto`.

---

## 11. Dependency updates

| Dependency | v3.7.4 | v4.1.x |
|---|---|---|
| `github.com/Masterminds/semver/v3` | v3.4.0 | v3.5.0 |
| `github.com/gin-contrib/logger` | v1.2.6 | v1.2.7 |
| `github.com/rs/zerolog` | v1.34.0 | v1.35.1 |
| `golang.org/x/net` | v0.52.0 | v0.57.0 |
| `golang.org/x/term` | v0.41.0 | v0.45.0 |
| `google.golang.org/grpc` | v1.79.3 | v1.82.1 |
| `github.com/quic-go/quic-go` | v0.59.0 | v0.59.1 |

---

## Quick migration checklist

- [ ] Update `go.mod` module path: `v3` → `v4`
- [ ] Update all import paths: `asset-link-sdk/v3/...` → `asset-link-sdk/v4/...`
- [ ] Replace `Identifiers(impl)` with `DeviceInfo(impl)` in the builder chain
- [ ] Implement `GetPropertyValues` and `GetSupportedProperties` instead of `GetIdentifiers`
  - Use `request.GetDevice().GetConnectionParameterSet()` to get connection parameters
  - Build a `DeviceInfo` from the retrieved device details
  - Call `deviceInfo.ConvertToPropertyValueResults()` to produce the response
- [ ] Update `registry.json`: replace `"siemens.common.identifiers.v1"` with `"siemens.connectivitysuite.deviceinfo.v1"` in the `app_types` array
- [ ] Update `model.NewDevice(...)` call sites to handle the returned `error`
- [ ] Update `AddNameplate(...)`: rename `manufacturerProductDesignation` argument to `productFamily`, handle returned `error`
- [ ] Update `AddNic(...)`: handle the returned `error` (always validates MAC format)
- [ ] Update `AddIPv4(...)`, `AddIPv6(...)`: handle the returned `error`; note that error fires only when **all** address fields are empty — passing an empty router address while IP and mask are set still succeeds
- [ ] Update `AddSoftware(...)` to `AddSoftwareArtifactComponent(...)`, handle returned `error`
- [ ] Use `AddProductInstanceIdentifier(vendor, articleNumber, serialNumber)` for product instance identifiers (replaces manual `ProductSerialIdentifier` construction)
- [ ] Update `AddDescription(...)`, `AddCapabilities(...)`: handle returned `error`
- [ ] Remove any references to `ManagementState`, `ReachabilityState`, `GatewayInfo`, `AssetContext`
- [ ] Remove any references to `ProductInstanceIdentifier`; use `ProductInstanceInformation` via `AddNameplate`
- [ ] Add optional `Description`, `DocUrl`, `FeedbackUrl` fields to `metadata.Metadata` (optional but recommended)
- [ ] Run `go mod tidy` and `go build ./...` to confirm no remaining compilation errors
