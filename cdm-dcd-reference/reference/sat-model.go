/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package reference

// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

import "fmt"
import "encoding/json"
import "reflect"

// Alphanumeric string that identifies a specific device family(type). For example:
// Maschinen Lesbare Fabrikatebezeichnung (MLFB).
type ArticleNumber string

type Capabilities struct {
  // boolean to determine if backup capability is supported.
  Backup *bool `json:"backup,omitempty"`

  // boolean to determine if change mode capability is supported.
  ChangeMode *bool `json:"change_mode,omitempty"`

  // If firmware update is supported
  FirmwareUpdate *bool `json:"firmware_update,omitempty"`

  // boolean to determine if memory reset capability is supported.
  MemoryReset *bool `json:"memory_reset,omitempty"`

  // boolean to determine if program update capability is supported.
  ProgramUpdate *bool `json:"program_update,omitempty"`

  // boolean to determine if reset communication parameter is supported.
  ResetCommunicationParameter *bool `json:"reset_communication_parameter,omitempty"`

  // boolean to determine if reset to factory capability is supported.
  ResetToFactory *bool `json:"reset_to_factory,omitempty"`

  // boolean to determine if restore capability is supported.
  Restore *bool `json:"restore,omitempty"`

  // TwoStepsFirmwareUpdate corresponds to the JSON schema field
  // "two_steps_firmware_update".
  TwoStepsFirmwareUpdate *bool `json:"two_steps_firmware_update,omitempty"`
}

// Description of the device.
type Description string

// Alphanumeric string that helps differentiate between different hardware versions
// of the same device model.
type DeviceHwVersion string

// Parameters that can be parsed by the discovery agent to create a thing
// description
type DeviceInfo struct {
  // ArticleNumber corresponds to the JSON schema field "article_number".
  ArticleNumber ArticleNumber `json:"article_number"`

  // Capabilities corresponds to the JSON schema field "capabilities".
  Capabilities *Capabilities `json:"capabilities,omitempty"`

  // Description corresponds to the JSON schema field "description".
  Description *Description `json:"description,omitempty"`

  // Detailed information about the device family
  DeviceDescription string `json:"device_description"`

  // Device-specific classification into families.
  DeviceFamily string `json:"device_family"`

  // DeviceHwVersion corresponds to the JSON schema field "device_hw_version".
  DeviceHwVersion DeviceHwVersion `json:"device_hw_version"`

  // DeviceSwVersion corresponds to the JSON schema field "device_sw_version".
  DeviceSwVersion DeviceSwVersion `json:"device_sw_version"`

  // Modules corresponds to the JSON schema field "modules".
  Modules []Module `json:"modules,omitempty"`

  // Nics corresponds to the JSON schema field "nics".
  Nics []DeviceInfoNicsElem `json:"nics,omitempty"`

  // PasswordProtected corresponds to the JSON schema field "password_protected".
  PasswordProtected PasswordProtected `json:"password_protected"`

  // Properties corresponds to the JSON schema field "properties".
  Properties *Properties `json:"properties,omitempty"`

  // SerialNumber corresponds to the JSON schema field "serial_number".
  SerialNumber SerialNumber `json:"serial_number"`

  // Alphanumeric string that identifies the producer or seller of the device.
  Vendor string `json:"vendor"`
}

type DeviceInfoNicsElem struct {
  // Device MAC address
  MacAddress *string `json:"mac_address,omitempty"`

  // NicIdentifier corresponds to the JSON schema field "nic_identifier".
  NicIdentifier *NicIdentifier `json:"nic_identifier,omitempty"`
}

// Alphanumeric string that identifies the Firmware running on the device.
type DeviceSwVersion string

type Module struct {
  // ArticleNumber corresponds to the JSON schema field "article_number".
  ArticleNumber ArticleNumber `json:"article_number"`

  // Capabilities corresponds to the JSON schema field "capabilities".
  Capabilities Capabilities `json:"capabilities"`

  // Description corresponds to the JSON schema field "description".
  Description Description `json:"description"`

  // DeviceHwVersion corresponds to the JSON schema field "device_hw_version".
  DeviceHwVersion DeviceHwVersion `json:"device_hw_version"`

  // DeviceSwVersion corresponds to the JSON schema field "device_sw_version".
  DeviceSwVersion DeviceSwVersion `json:"device_sw_version"`

  // Modules corresponds to the JSON schema field "modules".
  Modules []Module `json:"modules,omitempty"`

  // Name of the module
  Name string `json:"name"`

  // SerialNumber corresponds to the JSON schema field "serial_number".
  SerialNumber SerialNumber `json:"serial_number"`

  // Slot corresponds to the JSON schema field "slot".
  Slot Slot `json:"slot"`

  // StationNumber corresponds to the JSON schema field "station_number".
  StationNumber StationNumber `json:"station_number"`

  // SubSlot corresponds to the JSON schema field "sub_slot".
  SubSlot SubSlot `json:"sub_slot"`
}

// unique identifier of the network interface
type NicIdentifier string

// If device is password protected.
type PasswordProtected bool

type Properties struct {
  // ConnectedTo corresponds to the JSON schema field "connected_to".
  ConnectedTo []PropertiesConnectedToElem `json:"connected_to,omitempty"`

  // Connectivity status of the device.
  ConnectivityStatus PropertiesConnectivityStatus `json:"connectivity_status"`

  // Connectivity status of the device.
  DeviceType PropertiesDeviceType `json:"device_type"`

  // IpConfigurations corresponds to the JSON schema field "ip_configurations".
  IpConfigurations []PropertiesIpConfigurationsElem `json:"ip_configurations,omitempty"`

  // Operating mode of the device.
  OperatingMode *PropertiesOperatingMode `json:"operating_mode,omitempty"`

  // profinet name converted as per the requirment of the device.
  ProfinetName *string `json:"profinet_name,omitempty"`

  // RuntimeMode corresponds to the JSON schema field "runtime_mode".
  RuntimeMode *PropertiesRuntimeMode `json:"runtime_mode,omitempty"`

  // Slot corresponds to the JSON schema field "slot".
  Slot *Slot `json:"slot,omitempty"`

  // The Station Name of the device
  StationName string `json:"station_name"`

  // StationNumber corresponds to the JSON schema field "station_number".
  StationNumber *StationNumber `json:"station_number,omitempty"`

  // SubSlot corresponds to the JSON schema field "sub_slot".
  SubSlot *SubSlot `json:"sub_slot,omitempty"`

  // The version of the TIA from which the device was configured.
  TiapVersion *string `json:"tiap_version,omitempty"`
}

// Connection types
type PropertiesConnectedToElem struct {
  // Devices corresponds to the JSON schema field "devices".
  Devices []string `json:"devices"`

  // InterfaceType corresponds to the JSON schema field "interface_type".
  InterfaceType PropertiesConnectedToElemInterfaceType `json:"interface_type"`

  // Name corresponds to the JSON schema field "name".
  Name string `json:"name"`
}

type PropertiesConnectedToElemInterfaceType string

const PropertiesConnectedToElemInterfaceTypeEthernet PropertiesConnectedToElemInterfaceType = "ethernet"
const PropertiesConnectedToElemInterfaceTypeProfibus PropertiesConnectedToElemInterfaceType = "profibus"
const PropertiesConnectedToElemInterfaceTypeProfinet PropertiesConnectedToElemInterfaceType = "profinet"
const PropertiesConnectedToElemInterfaceTypeUnknown PropertiesConnectedToElemInterfaceType = "unknown"

type PropertiesConnectivityStatus string

const PropertiesConnectivityStatusOffline PropertiesConnectivityStatus = "offline"
const PropertiesConnectivityStatusOnline PropertiesConnectivityStatus = "online"
const PropertiesConnectivityStatusUnknown PropertiesConnectivityStatus = "unknown"

type PropertiesDeviceType string

const PropertiesDeviceTypeIndirectlyManaged PropertiesDeviceType = "indirectly_managed"
const PropertiesDeviceTypeNative PropertiesDeviceType = "native"
const PropertiesDeviceTypeNonManaged PropertiesDeviceType = "non-managed"

type PropertiesIpConfigurationsElem struct {
  // IpSettings corresponds to the JSON schema field "ip_settings".
  IpSettings []PropertiesIpConfigurationsElemIpSettingsElem `json:"ip_settings,omitempty"`

  // NicIdentifier corresponds to the JSON schema field "nic_identifier".
  NicIdentifier *NicIdentifier `json:"nic_identifier,omitempty"`
}

type PropertiesIpConfigurationsElemIpSettingsElem struct {
  // Default Network Gateway
  DefaultGateway string `json:"default_gateway"`

  // Device IP address
  IpAddress string `json:"ip_address"`

  // Device Subnet Mask
  SubnetMask string `json:"subnet_mask"`
}

type PropertiesOperatingMode string

const PropertiesOperatingModeNotsupported PropertiesOperatingMode = "notsupported"
const PropertiesOperatingModeRun PropertiesOperatingMode = "run"
const PropertiesOperatingModeStop PropertiesOperatingMode = "stop"

type PropertiesRuntimeMode string

const PropertiesRuntimeModeInMaintenance PropertiesRuntimeMode = "in_maintenance"
const PropertiesRuntimeModeNormal PropertiesRuntimeMode = "normal"

// Alphanumeric string that identifies a specific device instance within a specific
// device model.
type SerialNumber string

// Slot in which the device/module is present.
type Slot int

// The station number of the device
type StationNumber int

// Subslot number within the slot.
type SubSlot int

var enumValues_PropertiesConnectedToElemInterfaceType = []interface{}{
  "unknown",
  "ethernet",
  "profinet",
  "profibus",
}
var enumValues_PropertiesConnectivityStatus = []interface{}{
  "unknown",
  "online",
  "offline",
}
var enumValues_PropertiesDeviceType = []interface{}{
  "indirectly_managed",
  "non-managed",
  "native",
}
var enumValues_PropertiesOperatingMode = []interface{}{
  "notsupported",
  "run",
  "stop",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesOperatingMode) UnmarshalJSON(b []byte) error {
  var v string
  if err := json.Unmarshal(b, &v); err != nil {
    return err
  }
  var ok bool
  for _, expected := range enumValues_PropertiesOperatingMode {
    if reflect.DeepEqual(v, expected) {
      ok = true
      break
    }
  }
  if !ok {
    return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_PropertiesOperatingMode, v)
  }
  *j = PropertiesOperatingMode(v)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Module) UnmarshalJSON(b []byte) error {
  var raw map[string]interface{}
  if err := json.Unmarshal(b, &raw); err != nil {
    return err
  }
  if v, ok := raw["article_number"]; !ok || v == nil {
    return fmt.Errorf("field article_number: required")
  }
  if v, ok := raw["capabilities"]; !ok || v == nil {
    return fmt.Errorf("field capabilities: required")
  }
  if v, ok := raw["description"]; !ok || v == nil {
    return fmt.Errorf("field description: required")
  }
  if v, ok := raw["device_hw_version"]; !ok || v == nil {
    return fmt.Errorf("field device_hw_version: required")
  }
  if v, ok := raw["device_sw_version"]; !ok || v == nil {
    return fmt.Errorf("field device_sw_version: required")
  }
  if v, ok := raw["name"]; !ok || v == nil {
    return fmt.Errorf("field name: required")
  }
  if v, ok := raw["serial_number"]; !ok || v == nil {
    return fmt.Errorf("field serial_number: required")
  }
  if v, ok := raw["slot"]; !ok || v == nil {
    return fmt.Errorf("field slot: required")
  }
  if v, ok := raw["station_number"]; !ok || v == nil {
    return fmt.Errorf("field station_number: required")
  }
  if v, ok := raw["sub_slot"]; !ok || v == nil {
    return fmt.Errorf("field sub_slot: required")
  }
  type Plain Module
  var plain Plain
  if err := json.Unmarshal(b, &plain); err != nil {
    return err
  }
  *j = Module(plain)
  return nil
}

var enumValues_PropertiesRuntimeMode = []interface{}{
  "normal",
  "in_maintenance",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesRuntimeMode) UnmarshalJSON(b []byte) error {
  var v string
  if err := json.Unmarshal(b, &v); err != nil {
    return err
  }
  var ok bool
  for _, expected := range enumValues_PropertiesRuntimeMode {
    if reflect.DeepEqual(v, expected) {
      ok = true
      break
    }
  }
  if !ok {
    return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_PropertiesRuntimeMode, v)
  }
  *j = PropertiesRuntimeMode(v)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesConnectedToElemInterfaceType) UnmarshalJSON(b []byte) error {
  var v string
  if err := json.Unmarshal(b, &v); err != nil {
    return err
  }
  var ok bool
  for _, expected := range enumValues_PropertiesConnectedToElemInterfaceType {
    if reflect.DeepEqual(v, expected) {
      ok = true
      break
    }
  }
  if !ok {
    return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_PropertiesConnectedToElemInterfaceType, v)
  }
  *j = PropertiesConnectedToElemInterfaceType(v)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesIpConfigurationsElemIpSettingsElem) UnmarshalJSON(b []byte) error {
  var raw map[string]interface{}
  if err := json.Unmarshal(b, &raw); err != nil {
    return err
  }
  if v, ok := raw["default_gateway"]; !ok || v == nil {
    return fmt.Errorf("field default_gateway: required")
  }
  if v, ok := raw["ip_address"]; !ok || v == nil {
    return fmt.Errorf("field ip_address: required")
  }
  if v, ok := raw["subnet_mask"]; !ok || v == nil {
    return fmt.Errorf("field subnet_mask: required")
  }
  type Plain PropertiesIpConfigurationsElemIpSettingsElem
  var plain Plain
  if err := json.Unmarshal(b, &plain); err != nil {
    return err
  }
  *j = PropertiesIpConfigurationsElemIpSettingsElem(plain)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesConnectedToElem) UnmarshalJSON(b []byte) error {
  var raw map[string]interface{}
  if err := json.Unmarshal(b, &raw); err != nil {
    return err
  }
  if v, ok := raw["devices"]; !ok || v == nil {
    return fmt.Errorf("field devices: required")
  }
  if v, ok := raw["interface_type"]; !ok || v == nil {
    return fmt.Errorf("field interface_type: required")
  }
  if v, ok := raw["name"]; !ok || v == nil {
    return fmt.Errorf("field name: required")
  }
  type Plain PropertiesConnectedToElem
  var plain Plain
  if err := json.Unmarshal(b, &plain); err != nil {
    return err
  }
  *j = PropertiesConnectedToElem(plain)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Properties) UnmarshalJSON(b []byte) error {
  var raw map[string]interface{}
  if err := json.Unmarshal(b, &raw); err != nil {
    return err
  }
  if v, ok := raw["connectivity_status"]; !ok || v == nil {
    return fmt.Errorf("field connectivity_status: required")
  }
  if v, ok := raw["device_type"]; !ok || v == nil {
    return fmt.Errorf("field device_type: required")
  }
  if v, ok := raw["station_name"]; !ok || v == nil {
    return fmt.Errorf("field station_name: required")
  }
  type Plain Properties
  var plain Plain
  if err := json.Unmarshal(b, &plain); err != nil {
    return err
  }
  *j = Properties(plain)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesDeviceType) UnmarshalJSON(b []byte) error {
  var v string
  if err := json.Unmarshal(b, &v); err != nil {
    return err
  }
  var ok bool
  for _, expected := range enumValues_PropertiesDeviceType {
    if reflect.DeepEqual(v, expected) {
      ok = true
      break
    }
  }
  if !ok {
    return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_PropertiesDeviceType, v)
  }
  *j = PropertiesDeviceType(v)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PropertiesConnectivityStatus) UnmarshalJSON(b []byte) error {
  var v string
  if err := json.Unmarshal(b, &v); err != nil {
    return err
  }
  var ok bool
  for _, expected := range enumValues_PropertiesConnectivityStatus {
    if reflect.DeepEqual(v, expected) {
      ok = true
      break
    }
  }
  if !ok {
    return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_PropertiesConnectivityStatus, v)
  }
  *j = PropertiesConnectivityStatus(v)
  return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceInfo) UnmarshalJSON(b []byte) error {
  var raw map[string]interface{}
  if err := json.Unmarshal(b, &raw); err != nil {
    return err
  }
  if v, ok := raw["article_number"]; !ok || v == nil {
    return fmt.Errorf("field article_number: required")
  }
  if v, ok := raw["device_description"]; !ok || v == nil {
    return fmt.Errorf("field device_description: required")
  }
  if v, ok := raw["device_family"]; !ok || v == nil {
    return fmt.Errorf("field device_family: required")
  }
  if v, ok := raw["device_hw_version"]; !ok || v == nil {
    return fmt.Errorf("field device_hw_version: required")
  }
  if v, ok := raw["device_sw_version"]; !ok || v == nil {
    return fmt.Errorf("field device_sw_version: required")
  }
  if v, ok := raw["password_protected"]; !ok || v == nil {
    return fmt.Errorf("field password_protected: required")
  }
  if v, ok := raw["serial_number"]; !ok || v == nil {
    return fmt.Errorf("field serial_number: required")
  }
  if v, ok := raw["vendor"]; !ok || v == nil {
    return fmt.Errorf("field vendor: required")
  }
  type Plain DeviceInfo
  var plain Plain
  if err := json.Unmarshal(b, &plain); err != nil {
    return err
  }
  *j = DeviceInfo(plain)
  return nil
}
