---
title: "Artifact Management"
nav_order: 7
---

# Artifact management

## Artifact updates

To let your Asset Link apply artifacts such as firmware, software or device
configuration files, implement the three operations below. The gateway calls
Prepare first, then Activate to apply the artifact or Cancel to discard the
preparation.

| Operation | What your Asset Link must do |
| --- | --- |
| **Prepare** | Download and validate the update files, and prepare the device without installing the update or disrupting operation. |
| **Activate** | Install the prepared update and verify that it succeeded. |
| **Cancel** | If supported, discard the prepared update and clean up temporary files and device-side preparations. |

Want to try these operations without implementing a handler first? Follow the
[simulator quickstart](cmd.md#quickstart-update-a-simulated-device).

### Prerequisite: discover the target device

During discovery, populate each device's `DiscoveredDevice.ConnectionParameterSet`
with the schema URI and connection parameters your Asset Link needs to address
it. The gateway stores these parameters in Inventory and passes them back to
your artifact-update handler.

Use `GetConnectionParameterSet().GetParameterJson()` to identify and connect to
the target device, interpreting it in the same format your discovery
implementation produced. The SDK does not generate these parameters
automatically: `ConvertToDiscoveredDevice()` leaves the field empty, so set it
before publishing the discovered device. Without stored connection parameters,
the update service cannot create an update job for the device.

#### Connection parameter number limits

Your Asset Link defines the parameter format, for example `{"key":"device-1"}`.
If you use JSON numbers in connection or subdriver parameters, only values between
`-9007199254740991` and `9007199254740991` are supported.

### Implement the handler

If you use the [project template](bootstrap.md), select
`artifact_update=yes` to generate `handler/update.go`. Replace the stubs
with the download, validation and device communication logic for your devices.

For an existing Asset Link, add these methods to your handler:

```go
import "github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"

func (h *Handler) HandlePrepareUpdate(m *artifact.ArtifactMetaData, r artifact.ProgressReporter) (*artifact.Result, error)
func (h *Handler) HandleActivateUpdate(m *artifact.ArtifactMetaData, r artifact.ProgressReporter) (*artifact.Result, error)
func (h *Handler) HandleCancelUpdate(m *artifact.ArtifactMetaData, r artifact.ProgressReporter) (*artifact.Result, error)
```

Register your handler with `Update(handler)` to make artifact updates available
to the gateway:

```go
al := assetlink.New(metadata).
    Discovery(handler).
    Update(handler).
    Build()
```

When the Asset Link starts, the SDK automatically advertises
`siemens.industrialassethub.artifact_update.v1` so the gateway can find the
update service. No `registry.json` change is needed when using SDK registration.
If you register manually through `registry.json`, add the identifier to
`app_types`.

If the operation is not supported, return `artifact.Failed(artifact.NotSupported, "..."), nil` from that method.

### Working example with simulated devices

The [reference Asset Link][reference-main] registers its
[update handlers][reference-update] with the SDK. It
publishes each simulator address during discovery. `Prepare` downloads and
validates one simulated firmware file without changing the device. `Activate`
installs it, reboots the simulator and checks the active version. `Cancel`
discards the preparation; it does not roll back firmware.

A new job can reinstall the same firmware version. On retry, the simulator
only skips installation if it can confirm that this preparation, not just
the same version, was already installed.

The example accepts one prepared job at a time and keeps it in memory, so a
restart discards it. Run it in a test environment with a trusted OCI registry
and a writable temporary directory. Target unprotected simulators such as A0
or B0: the gateway update flow does not send device credentials. The
[reference tests][reference-tests] serve a firmware fixture from a local OCI
endpoint. They also show a lost device reply returning `Unknown` while the
gateway stream stays open. Simulated install and reboot stop when the gateway
context is cancelled. That does not undo an earlier step, and real devices
may continue working after the gateway call ends.

[reference-main]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/main.go
[reference-update]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/reference/update.go
[reference-tests]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/reference/update_test.go

### Request data

Each method receives `ArtifactMetaData` describing the target device, update job
and artifact (the update files). The SDK validates the request structure before
calling your handler. Your handler must still check that the device and update
content are supported.

| Getter | Content |
| --- | --- |
| `GetConnectionParameterSet()` | Device connection details provided during discovery: schema URI, JSON connection parameters and optional subdriver JSON. The current update flow does not supply device credentials. |
| `GetJobId()` | Update job ID. Use it to find the preparation associated with a later Activate or Cancel call. |
| `GetOciReference()` | Update location in an OCI registry: registry URL, repository and an immutable `sha256:` manifest digest. |

### Download the artifact

During Prepare, download the update files from the OCI registry. You can use the
`artifact/oci` package or another OCI client that verifies content sizes and
digests. The gateway sends the reference, not the files themselves.
The following snippet assumes `m` and `r` are your handler's metadata and
progress reporter, and `stagingDir` is a newly created private directory for this
Prepare attempt, separate from any previous successful preparation:

```go
import (
    "errors"
    "net/http"
    "os"

    "github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
)

ref := m.GetOciReference()
client := oci.Client{HTTP: http.DefaultClient}
files, err := client.Download(r.Context(), oci.Reference{
    RegistryURL:    ref.GetRegistryUrl(),
    Repository:     ref.GetRepository(),
    ManifestDigest: ref.GetManifestDigest(),
}, stagingDir)
if err != nil {
    err = errors.Join(err, os.RemoveAll(stagingDir))
    return artifact.Failed(artifact.ArtifactUnavailable, "Update files could not be downloaded"), err
}
// Validate the downloaded files for your device before returning success.
// files contains file names relative to stagingDir.
```

`Download` writes the payload files to `stagingDir` and verifies their sizes and
digests. A failed download can leave files behind; the example removes the whole
attempt directory and preserves any cleanup error alongside the download error.
Only use the files after the download succeeds. Your handler must also validate
the file format and compatibility with the target device, and clean up this
attempt's directory if validation or any later Prepare step fails. On success,
keep staged files if Activate or Cancel needs them.

The SDK downloader requires the OCI reference to point directly to an image
manifest describing the package's files. One manifest can describe multiple
files. OCI image indexes, which reference multiple manifests (for example,
platform-specific variants), are not supported. The downloader does not extract
archives; handle extraction yourself if your update format requires it.

### Report progress

```go
r.ReportIndeterminate("connecting")
r.ReportProgress("download", 42) // percent of this step, 0..100
```

- Step names are free text (max. 64 bytes); the percentage may restart per step.
- Progress at 100 % does **not** mean the operation finished; only the result does.
- Check the returned error, and stop work when `r.Context()` is done.

### Return the result

Return a result when the operation finishes. Choose the appropriate return
value for each path through your handler:

```go
return artifact.Succeeded(), nil
return artifact.Failed(artifact.ArtifactUnavailable, "Update files could not be downloaded"), err
return artifact.Unknown("Check the device before retrying"), err
```

| Result | Use it when |
| --- | --- |
| `Succeeded()` | The operation was verified as successful. For Prepare this means "ready to activate", not "installed". |
| `Failed(reason, msg)` | The operation did not start, or it failed **and** no device work is still running. Partial changes may remain. |
| `Unknown(msg)` | You cannot tell whether it worked or whether the device is still busy (for example connection lost during activation). |

Failure reasons: `InvalidRequest`, `NotSupported`, `TargetNotFound`,
`ArtifactUnavailable`, `ArtifactInvalid`, `DeviceRejected`, `OperationFailed`.

Always return an explicit result, even when you also return a Go error. Errors
are logged; they do not tell the gateway whether the update failed. A missing
or invalid result, a panic, or success accompanied by an error is reported as
`Unknown`.

Result messages (max. 512 bytes) are shown to users. Keep them actionable and
never include credentials, stack traces or raw device responses.

### Manage device state and cleanup

- Find the device using the connection parameters in every call.
- If Activate or Cancel needs something created by Prepare (for example
  staged files or a device-side upload), keep enough information to find it
  again. Associate the job ID with the target device and OCI reference, including
  the manifest digest, and reject later calls that conflict with that preparation.
  The SDK does not store this for you; decide whether it must survive a restart.
- Define how repeated Prepare, Activate and Cancel calls behave. Do not blindly
  repeat installation after an unknown outcome; check the device state first.
- Each gateway artifact-management-agent instance processes update operations
  sequentially. If other clients can call your Asset Link, or device-side work
  continues after a call ends, prevent conflicting operations on the same device.
- Clean up staged files and stored information after completion or, if supported,
  cancellation. An unsupported Cancel must return `Failed(NotSupported, ...)`
  without claiming that the preparation was removed.
- If the gateway aborts a call (`r.Context()` is done), that is **not** a
  cancel request. Only `HandleCancelUpdate` asks you to undo a preparation.

### Test your handler

Test each supported operation with a simulated device, covering success, failure
and loss of the device connection. Keep the gateway stream context active when
checking the result of a device connection loss; return `Unknown` if the device
outcome cannot be confirmed.

Test gateway-stream cancellation separately: make sure work stops when
`r.Context()` is done, without assuming that this undoes device-side work or
delivers a final result. `contracttest.Check` treats stream cancellation as a test
failure, so use a separate test for this behavior.

The `artifact/contracttest` helper runs your handler through the SDK wrapper and
checks the event stream and expected result, not the device behavior. Supply your
own simulated device and independently assert its state and interactions. Use a
context with a timeout; the helper runs handlers synchronously and cannot
interrupt a handler that ignores cancellation:

```go
import (
    "github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/contracttest"
    generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
)

contracttest.Check(t, ctx, request, handler.HandlePrepareUpdate,
    generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED)
```

Also check that Prepare does not install the update, Activate uses the correct
prepared files, and, if supported, Cancel removes the preparation without
installing it. For unsupported operations, check for `Failed(NotSupported, ...)`.
