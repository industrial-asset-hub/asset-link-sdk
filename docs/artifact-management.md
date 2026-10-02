---
title: "Artifact Management"
nav_order: 7
---

# Artifact management

## Artifact updates

For Gateway-managed artifact updates, your Asset Link must implement discovery
as well as the update handler. For every device you want to update, discovery
must report connection parameters that let your Asset Link reliably identify and
connect to that specific device. The Gateway stores those parameters with the
asset in Inventory and sends them back to your Asset Link when it runs an update.

The update handler applies artifacts such as firmware, software, or device
configuration files. The Gateway calls Prepare first, then Activate to apply
the artifact or Cancel to discard the preparation.

| Operation | What the update handler does |
| --- | --- |
| **Prepare** | Download and validate the update files from an OCI registry, and prepare the device without installing the update or disrupting operation. |
| **Activate** | Install the prepared update and verify that it succeeded. |
| **Cancel** | If supported, discard the prepared update and clean up temporary files and device-side preparations. |

`Cancel` only discards a prepared update before `Activate` starts. Once
`Activate` is in progress, let the update finish. Force-cancel is reserved for
stuck jobs; it does not guarantee that the device stops or rolls back the update.

To try the API before writing a handler, follow the
[simulator quickstart][simulator-quickstart], which uses firmware as an
example artifact.

## Quick path

Follow these steps to add artifact updates to a new Asset Link. Each step
links to the details if you need them.

1. **Generate the project.** Follow the [bootstrap guide](bootstrap.md) and
   select `artifact_update=yes`. The template registers the update handler and
   creates `handler/update.go` and `handler/update_test.go` with a working
   dummy implementation for the dummy device that discovery publishes. For an existing Asset Link, see
   [Add updates to an existing Asset Link](#add-updates-to-an-existing-asset-link).
2. **Publish connection parameters during discovery.** With
   `artifact_update=yes`, `createAndPublishMinimalDevice` in
   `handler/handler.go` already sets `ConnectionParameterSet` for the dummy
   device:

   ```go
   discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
   discoveredDevice.ConnectionParameterSet = &generated.ConnectionParameterSet{
       SchemaUri:     "urn:example:asset-link:mac-address:v1",
       ParameterJson: `{"macAddress":"00:16:3e:01:02:03"}`,
   }
   ```

   Set it the same way for each device your discovery publishes.
   Use your own schema URI and whatever your Asset Link needs to uniquely
   identify and reach the device. Details:
   [Connection parameters](#connection-parameters).
3. **Implement the handlers.** `handler/update.go` already works for the
   dummy device. Replace the parts marked with `TODO` with the logic for your
   devices, and adapt `connectionSchema` and `connectionParameters` to the
   format from step 2. The dummy implementation:
   - reads the connection parameters in every call,
   - downloads and verifies the artifact from the OCI registry in Prepare,
   - remembers the preparation per job ID in memory,
   - "installs" it on an in-memory dummy device in Activate, or discards it
     in Cancel, and cleans up.

   Details: [Request data](#request-data),
   [Download the artifact](#download-the-artifact),
   [Return the result](#return-the-result),
   [Manage device state and cleanup](#manage-device-state-and-cleanup).
4. **Update the unit tests.** `handler/update_test.go` serves a test artifact
   from a local HTTP server and checks Prepare, Activate, Cancel and an unknown
   device with `contracttest.Check`. Adapt it as you replace the dummy logic.
   Details: [Test your handler](#test-your-handler).
5. **Try it locally with `al-ctl`.** Start your Asset Link and follow
   [Test artifact updates on your Asset Link](cmd.md#test-artifact-updates-on-your-asset-link):
   start a registry, upload a file, and run Prepare and Activate.

For Gateway-managed updates, the Gateway uses the connection parameters stored
during discovery. If a device was discovered before you added them, discover
it again.

## Details

### Connection parameters

`SchemaUri` identifies your connection parameter format; it is not a built-in
SDK schema. In `ParameterJson`, report whatever your Asset Link needs to
uniquely identify and reach the device. A MAC address is just one example. If
the connection needs subdriver configuration, include it as a JSON string in
the optional `SubdriverJson` field. See the
[reference Asset Link's `Discover` method][reference-discovery] for an example.

In the template, `generated` is already imported in `handler/handler.go` as
`github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery`.
`ConvertToDiscoveredDevice()` leaves `ConnectionParameterSet` empty, so set it
before calling `PublishDevice`.

Every update call contains these same parameters. Without stored connection
parameters, the update service cannot create a job for the device.

If you use JSON numbers in connection or subdriver parameters, only values
between `-9007199254740991` and `9007199254740991` are supported.

### Add updates to an existing Asset Link

The easiest start is the dummy implementation from the template:

1. Publish connection parameters during discovery. See
   [Connection parameters](#connection-parameters).
2. Generate a project with `artifact_update=yes` (see the
   [bootstrap guide](bootstrap.md)) and copy `handler/update.go` and
   `handler/update_test.go` into your Asset Link.
3. Rename the receiver `*AssetLinkImplementation` to your handler type.
4. Adapt `connectionSchema` and `connectionParameters` to the format your
   discovery publishes.
5. Remove the `dummyDeviceMAC` check and the `dummyDevice` state, and add your
   device logic at the `TODO`s. Otherwise every real device fails with
   `TargetNotFound`, or the update only changes in-memory state.
6. Register the handler with `Update(handler)`, as shown below.

Your handler type needs these methods:

```go
import "github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"

func (h *Handler) HandlePrepareUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error)
func (h *Handler) HandleActivateUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error)
func (h *Handler) HandleCancelUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error)
```

Register the handler with `Update(handler)`:

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

If an operation is not supported, return
`artifact.Failed(artifact.NotSupported, "..."), nil` from that method.

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

### Publish an artifact to an OCI registry

Use [ORAS](https://oras.land/docs/installation/) to push the update files from
the directory that contains them. For the [local test registry][local-oci-registry]:

```bash
oras push --plain-http localhost:5000/my-update:test update.bin
```

Replace `update.bin` with your file and `my-update:test` with your repository
and tag.
For an HTTPS registry, use its host and omit `--plain-http`.

Record the registry URL (`http://localhost:5000` in this example), the
repository (`my-update`), and the `sha256:…` **manifest digest** printed by
ORAS. Use that digest in the update request, not the tag or a file checksum.

### Download the artifact

During Prepare, download the update files from the OCI registry. You can use the
`artifact/oci` package or another OCI client that verifies content sizes and
digests. The gateway sends the reference, not the files themselves.

The generated `handler/update.go` already contains this code. Otherwise, add these imports:

```go
import (
    "errors"
    "net/http"
    "os"

    "github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
)
```

Then add the download code inside `HandlePrepareUpdate`. It creates a fresh private directory
for each download attempt, separate from any previously prepared update:

```go
stagingDir, err := os.MkdirTemp("", "asset-link-update-")
if err != nil {
    return artifact.Failed(artifact.OperationFailed, "Could not create a staging directory"), err
}

ref := metadata.GetOciReference()
client := oci.Client{HTTP: http.DefaultClient}
files, err := client.Download(transmitter.Context(), oci.Reference{
    RegistryURL:    ref.GetRegistryUrl(),
    Repository:     ref.GetRepository(),
    ManifestDigest: ref.GetManifestDigest(),
}, stagingDir)
if err != nil {
    err = errors.Join(err, os.RemoveAll(stagingDir))
    return artifact.Failed(artifact.ArtifactUnavailable, "Update files could not be downloaded"), err
}
// Continue with device-specific validation and staging here.
// Use filepath.Join(stagingDir, name) to open each name in files.
```

`Download` verifies file sizes and digests, but your handler still needs to
validate the file format and device compatibility before returning success.
For example, return `Failed(ArtifactInvalid, ...)` if the artifact contains no
usable files.
Remove `stagingDir` if validation or any later Prepare step fails, just as the
snippet does for a failed download. If Activate or Cancel needs the files,
save the directory path with the prepared job and clean it up when the job
completes or is cancelled. Otherwise, remove it before Prepare returns.

The empty first argument to `os.MkdirTemp` uses the system temporary directory.
The reference Docker image and images generated with `artifact_update=yes`
include a writable `/tmp`. If you run
with a read-only root filesystem, mount a writable `/tmp` (for example,
`docker run --read-only --tmpfs /tmp:mode=1777 …`) or set `TMPDIR` to a writable
mounted directory.

If the preparation must survive a restart, pass a configured, writable persistent
directory instead and persist the job state that records the staging path.

The SDK downloader requires the OCI reference to point directly to an image
manifest describing the package's files. One manifest can describe multiple
files. OCI image indexes, which reference multiple manifests (for example,
platform-specific variants), are not supported. The downloader does not extract
archives; handle extraction yourself if your update format requires it.

### Report progress

```go
transmitter.ReportIndeterminate("connecting")
transmitter.ReportProgress("download", 42) // percent of this step, 0..100
```

- Step names are free text (max. 64 bytes); the percentage may restart per step.
- Progress at 100 % does **not** mean the operation finished; only the result does.
- Check the returned error, and stop work when `transmitter.Context()` is done.

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
- Once Activate has started, reject Cancel and a new Prepare for the same job:
  Cancel cannot undo changes Activate may already have made. Allow Activate to
  be retried instead. The generated `handler/update.go` does this.
- Each gateway artifact-management-agent instance processes update operations
  sequentially. If other clients can call your Asset Link, or device-side work
  continues after a call ends, prevent conflicting operations on the same device.
- Clean up staged files and stored information after completion or, if supported,
  cancellation. An unsupported Cancel must return `Failed(NotSupported, ...)`
  without claiming that the preparation was removed.
- If the gateway aborts a call (`transmitter.Context()` is done), that is **not** a
  cancel request. Only `HandleCancelUpdate` asks you to undo a preparation.

### Test your handler

The generated `handler/update_test.go` shows how to build a request and call
each handler. Its `serveUpdateArtifact` helper serves a test artifact from a
local OCI server; reuse it in your own tests. Test each supported operation with a simulated device, covering success, failure
and loss of the device connection. Keep the gateway stream context active when
checking the result of a device connection loss; return `Unknown` if the device
outcome cannot be confirmed.

Test gateway-stream cancellation separately: make sure work stops when
`transmitter.Context()` is done, without assuming that this undoes device-side work or
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
[reference-discovery]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/reference/reference.go#L116-L140
[reference-update]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/reference/update.go
[reference-tests]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/reference/update_test.go
[simulator-quickstart]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/README.md#update-firmware-on-a-simulated-device
[local-oci-registry]: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/cdm-al-reference/README.md#start-a-local-oci-registry
