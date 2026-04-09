---
title: "Firmware Update"
nav_order: 7
---

### Firmware Update

#### Introduction

- The firmware update is implemented as a two-step process:
   1. Prepare (initiates the installation process):
      Typically installs the new version of the firmware alongside the old one if the device supports a two-step update.
      This step requires some time but it can often be performed during regular operation of the device.

   2. Activate (finalizes the installation process and is typically triggered manually):
      Actually activates the new version (e.g., by switching the desired firmware version and rebooting the device).
      This activation step typically comes with a short interruption of the regular device operation and can be performed when the time is suitable.
      However, for legacy devices that do not support a two-step update the activate step might also perform the installation (as the update/firmware artefact is available during both steps) which might lead to a more lengthy operation.

    The update process can fail when the asset link reports an error and it can be explicitly canceled by the user after the prepare step (which would allow for a cleanup of the prepared update and installed files).

- Finding Devices via the Device Identifier Blob:
  Alongside the actual firmware update process the `feat/update` branch of the SDK also extends the discovered assets and adds a vendor-specific device identifier blob during discovery. This device identifier blob is provided to the asset link when an update is triggered and allows the asset link to find the relevant device. It can include vendor-specific data like device IDs, bus IDs, or connection parameters.

- Alongside firmware updates the `feat/update` branch of the SDK also adds the ability to transfer arbitrary artefacts (i.e., logs, backups, configurations) from/to the device

#### Overview of Simulated Devices (Reported/Supported by the Reference Asset Link)

| Device                     | Interface | Parent Device       | Discovery Credentials | Update Credentials |
|----------------------------|-----------|---------------------|-----------------------|--------------------|
| Simulated Device A0        | eth0      | -                   | not required          | not required       |
| Simulated Device A1        | eth0      | -                   | not required          | **required**       |
| Simulated Device A2        | eth0      | -                   | **required**          | **required**       |
| Simulated Device B0        | eth1      | -                   | not required          | not required       |
| Simulated Sub Device B0-0  | -         | Simulated Device B0 | n/a                   | n/a                |
| Simulated Sub Device B0-1  | -         | Simulated Device B0 | n/a                   | n/a                |
| Simulated Sub Device B0-2  | -         | Simulated Device B0 | n/a                   | n/a                |
| Simulated Device B1        | eth1      | -                   | not required          | **required**       |
| Simulated Device B2        | eth1      | -                   | **required**          | **required**       |

**Read Credentials (required for Discovery):**
- Username: `user`
- Passowrd: `user_passowrd`

**Admin Credentials (required for Update):**
- Username: `admin`
- Password: `admin_password`

#### Relevant Files
The firmware update files for the simulated devices and their PCN data for the AAS can be found here:
https://github.com/industrial-asset-hub/asset-link-sdk/tree/feat/update/misc

#### Local Testing with `alctl`
Local testing does not require additional infrastructure instead it utilizes the `alctl` tool (that was extended with `update` operations) to interact with the asset link (e.g., the reference asset link) directly.

##### Preparations
```bash
# Switch to the SDK folder

# Start the Reference Asset Link (in a dedicated terminal)
go run cdm-al-reference/main.go
```

##### Perform Firmware Update without Credentials
```bash
# Trigger the prepare step of the firmware update for a device (no device credentials required)
go run cmd/al-ctl/al-ctl.go update prepare --job-id 0 --device-identifier-file misc/device_address.json --convert-device-identifier --artefact-type firmware --artefact-file misc/simulated_device_firmware_2.0.0.fwu

# Trigger the activate step of the firmware update for a device (no device credentials required)
go run cmd/al-ctl/al-ctl.go update activate --job-id 0 --device-identifier-file misc/device_address.json --convert-device-identifier --artefact-type firmware --artefact-file misc/simulated_device_firmware_2.0.0.fwu

# Cancel an update after the prepare step (no device credentials required)
go run cmd/al-ctl/al-ctl.go update cancel --job-id 0 --device-identifier-file misc/device_address.json --convert-device-identifier --artefact-type firmware
```

##### Perform Firmware Update with Credentials
```bash
# Trigger the prepare step of the firmware update for a device (device credentials required)
go run cmd/al-ctl/al-ctl.go update prepare --job-id 1 --device-identifier-file misc/device_address_cred.json --convert-device-identifier --device-credentials-file misc/device_credentials_admin.json --artefact-type firmware --artefact-file misc/simulated_device_firmware_2.0.0.fwu

# Trigger the activate step of the firmware update for a device (device credentials required)
go run cmd/al-ctl/al-ctl.go update activate --job-id 1 --device-identifier-file misc/device_address_cred.json --convert-device-identifier --device-credentials-file misc/device_credentials_admin.json --artefact-type firmware --artefact-file misc/simulated_device_firmware_2.0.0.fwu

# Cancel an update after the prepare step (device credentials required)
go run cmd/al-ctl/al-ctl.go update cancel --job-id 1 --device-identifier-file misc/device_address_cred.json --convert-device-identifier --device-credentials-file misc/device_credentials_admin.json --artefact-type firmware
```

##### Remarks
- The `Simulated Device A2` and `Simulated Device B2` will not be discovered as no device (read) credentials are provided during discovery.
- Only regular simulated devices support an update (i.e., simulated sub-devices do not support an update).
- Only the `Simulated Device A1` and `Simulated Device B1` need device (admin) credentials (username: admin password: admin_password) for updates.
- During an update the new and old/current firmware version cannot be the same or the update will fail with a corresponding error.
