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
HandlePrepareUpdate(*artefact.ArtefactMetaData, artefact.StatusTransmitter) error
HandleActivateUpdate(*artefact.ArtefactMetaData, artefact.StatusTransmitter) error
HandleCancelUpdate(*artefact.ArtefactMetaData, artefact.StatusTransmitter) error
```

The service is registered and advertised to the registry only when an update
implementation is supplied. As with discovery, a custom generated
`ArtefactUpdateApiServer` can be assigned to the builder instead; it takes
precedence over the SDK adapter.

## Control API

All three RPCs take **one `ArtefactMetaData` request** and return a server stream
of `ArtefactOperationStatus`. There is no client stream or file-transfer protocol.

```proto
rpc PrepareUpdate(ArtefactMetaData) returns (stream ArtefactOperationStatus);
rpc ActivateUpdate(ArtefactMetaData) returns (stream ArtefactOperationStatus);
rpc CancelUpdate(ArtefactMetaData) returns (stream ArtefactOperationStatus);
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

All handlers report phase, state, message, and progress (0–100) using
`transmitter.UpdateStatus`. Check its returned error and use
`transmitter.Context()` to observe RPC cancellation and deadlines. Malformed
metadata is rejected with `InvalidArgument` before the handler is called.
Handler errors are returned through gRPC; use `status.Error` for an appropriate
error code. EOF ends a successful handler's status stream; clients must still
check the operation statuses rather than treating transport success as an update.

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
