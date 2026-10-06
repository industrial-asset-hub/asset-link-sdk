---
title: "Software Update"
nav_order: 7
---

# Software update

Asset links can implement the update RPCs of
`siemens.industrialassethub.artefact_update.v1.ArtefactUpdateApi`:

- **PrepareUpdate** prepares or installs software for later activation.
- **ActivateUpdate** activates a prepared update. Devices without two-stage
  updates can request the software here and perform installation at this point.
- **CancelUpdate** cleans up a prepared update.
- **GetUpdateTargetRequirements** optionally describes which fields of an
  Inventory v2 asset are needed to identify the update target.

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

## Update target requirements (optional)

An Asset Link can additionally implement
`artefact.UpdateTargetRequirementsProvider` on its update handler. The existing
three update-handler methods remain unchanged. When the optional interface is
absent, `GetUpdateTargetRequirements` returns gRPC `Unimplemented`. Its request
contains the artefact type (e.g. `AT_FIRMWARE`), not a device identifier.

The response declares an `identifier_format_uri` and a list of required fields.
Each `asset_json_pointer` is an RFC 6901 JSON Pointer relative to the `data`
object returned by Inventory `GET /v2/assets/{id}`. For example, SAT might
request `article_number` from
`/product_instance_information/manufacturer_product/product_id` and
`serial_number` from `/product_instance_information/serial_number`.
The cloud resolves each field to a non-empty string and constructs the
versioned, Asset-Link-specific device identifier. For the initial contract,
the decoded identifier bytes contain UTF-8 JSON with this envelope (the Asset
Link interprets the values):

```json
{
  "identifier_format_uri": "urn:siemens:iah:sat-driver:update-target:v1",
  "values": {"article_number": "6ES7...", "serial_number": "ABC123"}
}
```

The `values` keys correspond to the declared field names; every value is a
non-empty string. No device values are returned by the requirements RPC.
Properties not present in the Inventory
asset read representation (such as `asset_identifiers`) cannot be selected
through this contract. Do not use paths to credentials or secrets.

The agent can stage an update and publish these requirements before the cloud
populates the identifier. Since `PrepareUpdate` still requires device
identifier metadata, the actual device-side prepare and activate calls must
both occur after the identifier is available. The populated identifier is
carried in the existing `DeviceIdentifier.blob` for both calls; its payload
format is owned by the Asset Link. The current SDK expects the wire-level
`blob` to contain base64-encoded identifier bytes and supplies decoded bytes
to the handler via `GetDeviceIdentifierBlob()`.

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
