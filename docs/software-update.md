---
title: "Software Update"
nav_order: 7
---

# Software update

Asset links can implement the three update RPCs of
`siemens.industrialassethub.artefact_update.v1.ArtefactUpdateApi`:

- **PrepareUpdate** prepares or installs software for later activation.
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
HandlePrepareUpdate(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error
HandleActivateUpdate(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error
HandleCancelUpdate(artefact.ArtefactMetaData, artefact.StatusTransmitter) error
```

The service is registered and advertised to the registry only when an update
implementation is supplied. As with discovery, a custom generated
`ArtefactUpdateApiServer` can be assigned to the builder instead; it takes
precedence over the SDK adapter.

## Streaming protocol

Prepare and activate use bidirectional streams. The client first sends an
`ArtefactChunk` containing `ArtefactMetaData`, including job, artifact, and device
identifiers. Credential messages are optional.
`DeviceIdentifier.blob` contains base64-encoded device-specific bytes;
the SDK decodes them for `GetDeviceIdentifierBlob()`.

When a handler calls `ReceiveArtefactToWriter`, `ReceiveArtefactToFile`, or
`ReceiveArtefactToData`, the SDK sends an `AORT_ARTEFACT_TRANSMISSION` request.
The client then sends file-content chunks and closes its send side to indicate
the end of the file. The server can continue sending status messages until the
handler returns. Receiving the artifact is optional in either stage. Prefer
the writer/file helpers for large artifacts, since the data helper buffers the
entire file in memory.

Cancel takes metadata in a single request and returns a stream of operation
statuses. All handlers can report phase, state, message, and progress (0–100)
using `UpdateStatus`. Check its returned error and use `Context()` to observe
RPC cancellation and deadlines.

Malformed metadata or unexpected content chunks are rejected with
`InvalidArgument`. Handler errors are returned through gRPC, so handlers can
use `status.Error` to choose an appropriate error code.

The asset-link implementation owns job persistence, device-specific validation,
serialization of updates to the same device, and cleanup. A CancelUpdate RPC
is a device-level operation, distinct from cancelling an RPC context.
