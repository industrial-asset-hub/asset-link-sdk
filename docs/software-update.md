---
title: "Software Update"
nav_order: 7
---

# Software update

Asset links can implement the update RPCs of
`siemens.industrialassethub.artefact_update.v1.ArtefactUpdateApi`:

- **PrepareUpdate** performs non-disruptive preparation for later activation.
- **ActivateUpdate** activates a prepared update. Devices without two-stage
  updates can request the software here and perform installation at this point.
- **CancelUpdate** cleans up a prepared update.

Register an implementation with the builder:

```go
al := assetlink.New(metadata).
    Discovery(handler).
    Update(handler).
    Build()
```

The handler implements these methods using the public `artefact` package:

```go
HandlePrepareUpdate(*artefact.ArtefactMetaData, artefact.ProgressReporter) (*artefact.Result, error)
HandleActivateUpdate(*artefact.ArtefactMetaData, artefact.ProgressReporter) (*artefact.Result, error)
HandleCancelUpdate(*artefact.ArtefactMetaData, artefact.ProgressReporter) (*artefact.Result, error)
```

The service is registered and advertised to the registry only when an update
implementation is supplied. As with discovery, a custom generated
`ArtefactUpdateApiServer` can be assigned to the builder instead; it takes
precedence over the SDK adapter.

## Control API

All three RPCs take **one `ArtefactMetaData` request** and return `UpdateEvent`
streams using the same SDK wrapper and validation rules. Device-level cancellation
is not invoked by the gateway's prepared-job cancellation workflow. There is no
client stream or file-transfer protocol.

```proto
rpc PrepareUpdate(ArtefactMetaData) returns (stream UpdateEvent);
rpc ActivateUpdate(ArtefactMetaData) returns (stream UpdateEvent);
rpc CancelUpdate(ArtefactMetaData) returns (stream UpdateEvent);
```

Metadata has six direct fields:

| Field | Purpose |
| --- | --- |
| `job_id` | Required operation identity. |
| `asset_id`, `view_id` | Optional Inventory context for traceability. |
| `connection_parameter_set` | Required typed Connectivity Suite target: schema URI, parameter JSON, optional subdriver JSON and connector credentials. |
| `artefact_type` | Firmware, software, etc. |
| `oci_reference` | Required local registry URL, repository and immutable sha256 manifest digest. |

Handlers receive the same typed metadata as the gRPC request, without a
conversion layer. Connection parameters are available through
`metadata.GetConnectionParameterSet()` and the artifact through
`metadata.GetOciReference()`. There are no identifier blobs, base64 encoding or
duplicate credential messages. Upstream registry credentials belong to the
local OCI cache configuration, not the update request.

## Shared update result contract

Handlers report step progress through `reporter.ReportProgress("download", 42)`
or `reporter.ReportIndeterminate("connecting")`. Percentages are per step and may
reset between steps. Check reporting errors and observe `reporter.Context()` for
deadlines. Progress at 100 percent is **not** completion.

Return one typed result; the SDK owns final result emission:

```go
return artefact.Succeeded(), nil
return artefact.Failed(artefact.ArtifactUnavailable, "Firmware could not be retrieved"), err
return artefact.Unknown("Verify the device before another update"), err
```

| Outcome | Contract |
| --- | --- |
| `SUCCEEDED` | Success of the requested operation was verified. Prepare success means ready for activation, not installed. |
| `FAILED` | The requested change did not start, or failure and the absence of ongoing device work were confirmed. Partial changes may remain; no rollback is promised. |
| `UNKNOWN` | Completion and/or ongoing device work cannot be confirmed. No automatic replay is safe. |

Nil/unspecified/invalid results, ordinary errors without a confirmed failure, and
handler panics default to `UNKNOWN`. Diagnostic errors are logged, not copied into
public messages. A success result accompanied by an error becomes unknown. An
explicit failed result may include an error for diagnostic logging.

Outcomes describe the requested operation. A failed/unsupported cancellation does
not prove that a preceding activation stopped. Consumers must require a successful
cancellation before treating it as completed device cleanup. Unsupported device
cancellation returns `artefact.Failed(artefact.NotSupported, ...)`, not success.

The wire stream contains zero or more progress events followed by exactly one
result and normal EOF. Empty events, invalid enums/percentages, missing or duplicate
results and events after the result violate the contract. A result followed by a
transport error is not authoritative. `artefact.ReadResult` validates this protocol
for clients and custom servers; custom servers can use `artefact.RunUpdate` to get
the same conservative wrapper. Failure to send progress/results remains a stream
error, even if a handler ignores it. Reporters cannot send after handler completion.

Step names are informational (maximum 64 bytes). Public failure messages are
optional and limited to 512 UTF-8 bytes. Never include credentials, response bodies
or stack traces. Failure reasons are structured codes; clients must not parse
messages or generic gRPC errors to infer the device outcome.

Malformed metadata is rejected before the handler is called. RPC cancellation
does not prove that device work stopped and is not a device cancellation request.

## Customer contract tests

The reusable `artefact/contracttest.Check` harness executes a fixture-backed handler
through the real SDK wrapper, validates its stream and checks the expected outcome:

```go
contracttest.Check(t, ctx, request, handler.HandlePrepareUpdate,
    generated.UpdateOutcome_UPDATE_OUTCOME_FAILED)
```

Use independent success, missing-artifact/rejection and lost-response fixtures with
a bounded context. Verify device calls in those fixtures: a protocol test cannot
prove that a driver telling the SDK "success" actually waited for the device.

This is a breaking change to the unreleased update API. Regenerate clients from
`specs/artefact_update.proto`, migrate all three handlers and update SDK pins
together. Discovery and DeviceInfo APIs are unchanged.

## OCI downloads

**No firmware bytes are sent over gRPC.** Downloading and staging are Asset Link
decisions. Prepare must not install firmware or disrupt the device. The
`artefact/oci` helper resolves and verifies single OCI manifests, checks complete
blob availability and streams payloads to files with bounded memory. Downloads
retry up to three times from the beginning; resumable range requests, image
indexes and archive extraction are not implemented. Use a private staging
directory and publish it only after successful verification of all payloads.

The asset-link implementation owns job persistence, device-specific validation,
serialization of updates to the same device, and cleanup. A CancelUpdate RPC
is a device-level operation, distinct from cancelling an RPC context.
