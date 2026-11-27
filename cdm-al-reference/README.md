# Reference Asset Link

This example Asset Link discovers simulated devices and demonstrates the
[artifact update API](../docs/artifact-management.md) using sample firmware.
For testing updates on your own Asset Link, see the
[generic `al-ctl` workflow](../docs/cmd.md#test-artifact-updates-on-your-asset-link).

## Update firmware on a simulated device

This walkthrough updates simulated device **A0** from firmware `1.0.0` to
`2.0.0` on a fresh simulator. No physical device or full Asset Gateway is
needed. Use a local SDK checkout, the Go toolchain required by the repo,
Docker, and [ORAS](https://oras.land/docs/installation/). Make sure the
[`al-ctl` executable](../docs/cmd.md#command-line-tool-for-local-debugging)
is on your `PATH`. Run the commands from the repository root
(`asset-link-sdk/`) in both terminals.

### 1. Start the demo (terminal 1)

#### Start a local OCI registry

For local testing, use [`docker.io/library/registry:3`](https://hub.docker.com/_/registry),
the Docker Official Image for the Distribution registry. It stores the sample
firmware on port `5000`. Skip `docker run` if a registry is already running on
that port.

This command publishes the port on `127.0.0.1`. The registry uses plain HTTP
without authentication, so other processes on the host can read or write
artifacts. No persistent storage is configured; treat uploads as disposable.
Use only non-sensitive test artifacts. For shared or production use, choose
an approved OCI registry with TLS and authentication.

The `:3` tag can change. If your policy requires an exact, reviewed image,
use an approved image digest instead of the tag.

```bash
docker run --rm -d --name asset-link-demo-oci \
  -p 127.0.0.1:5000:5000 docker.io/library/registry:3
```

#### Start the reference Asset Link

```bash
go run ./cdm-al-reference
```

Leave this terminal running. Open the [simulator UI](http://localhost:8083/)
and click **+** next to **Simulated Device A0**. Before continuing, check that
**Installed Firmware Version** and **Active Firmware Version** both show
`1.0.0` and that **Device State** is `active`. If you reused a simulator that
has already been updated, restart the reference Asset Link to reset the
simulated devices.

### 2. Publish the firmware and fill in the example request (terminal 2)

Run the push from `misc/` as shown so that the OCI layer name contains only
the filename, not a path.

```bash
digest=$(
  cd misc &&
  oras push --plain-http localhost:5000/reference-firmware:2.0.0 \
    simulated_device_firmware_2.0.0.fwu | awk '/^Digest:/ {print $2}'
)
printf 'Manifest digest: %s\n' "$digest"
```

If the printed digest is empty, check that the OCI registry is running and
the ORAS push succeeded before continuing.

Use the [example update request](../misc/artifact_update_request.json) to
create a temporary request file with a new `job_id` and the `manifest_digest`
from the push. Keep the simulator's `schema_uri` and its `parameter_json`
string as shown in the example. This leaves the checked-in example unchanged.

```bash
request=$(mktemp)
sed -e "s|JOB_ID|reference-update-$(date +%s)|" \
  -e "s|MANIFEST_DIGEST|${digest:?Run the ORAS push first}|" \
  misc/artifact_update_request.json > "$request"
```

### 3. Prepare, then activate or cancel (terminal 2)

```bash
al-ctl update prepare --endpoint localhost:8081 --request "$request" --timeout 20
```

Look for `Progress: "download"` and `SUCCEEDED`. Both firmware fields for A0
should still show `1.0.0`. After Prepare succeeds, use the **same request
file** to activate:

```bash
al-ctl update activate --endpoint localhost:8081 --request "$request" --timeout 20
```

Watch A0 in the UI: it should show `updating` and `booting` before returning
to `active`. Both firmware fields should now show `2.0.0`. If you want to
test Cancel instead, run this **in place of Activate**:

```bash
al-ctl update cancel --endpoint localhost:8081 --request "$request" --timeout 20
```

Cancel should return `SUCCEEDED` without changing the firmware version. To
repeat the demo, create a request with a new job ID. A new job can reinstall
the same version, so the version number may stay the same.

Progress and the result appear in **terminal 2**. `SUCCEEDED` and exit status 0
confirm success. `FAILED` and `UNKNOWN` return a nonzero status. If the outcome
is unknown, check the device before retrying.

### 4. Clean up

Press `Ctrl+C` in terminal 1. If you started the OCI registry for this
walkthrough, stop it with `docker stop asset-link-demo-oci`.
