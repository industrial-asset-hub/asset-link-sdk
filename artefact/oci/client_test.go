// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package oci

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveEnsureAndDownload(t *testing.T) {
	payload, config := []byte("firmware bytes"), []byte("{}")
	m := Manifest{SchemaVersion: 2, MediaType: ManifestMediaType, Config: Descriptor{Digest: digestOf(config), Size: int64(len(config))}, Layers: []Descriptor{{Digest: digestOf(payload), Size: int64(len(payload)), Annotations: map[string]string{"org.opencontainers.image.title": "firmware.upd"}}}}
	manifest, _ := json.Marshal(m)
	blobs := map[string][]byte{digestOf(config): config, digestOf(payload): payload}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/manifests/") {
			w.Write(manifest)
			return
		}
		w.Write(blobs[filepath.Base(r.URL.Path)])
	}))
	defer server.Close()
	c := Client{HTTP: server.Client()}
	r, got, err := c.Resolve(context.Background(), server.URL, "firmware/s7-1500", "latest")
	if err != nil || r.ManifestDigest != digestOf(manifest) {
		t.Fatalf("resolve: %+v %v", r, err)
	}
	if err := c.EnsureAvailable(context.Background(), r, got); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	files, err := c.Download(context.Background(), r, dir)
	if err != nil || len(files) != 1 {
		t.Fatalf("download: %v %v", files, err)
	}
	data, err := os.ReadFile(filepath.Join(dir, files[0]))
	if err != nil || string(data) != string(payload) {
		t.Fatalf("payload: %q %v", data, err)
	}
}
func TestRejectUnsafeReferencesAndPayloadNames(t *testing.T) {
	for _, registry := range []string{"file:///tmp", "http://user:password@localhost", "http://localhost/path", "http://localhost?x=1"} {
		if err := (Reference{RegistryURL: registry, Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}).Validate(); err == nil {
			t.Fatalf("accepted %s", registry)
		}
	}
	for _, name := range []string{"../escape", "sub/file", `sub\file`, "..", "prepared.json"} {
		if _, err := payloadName(Descriptor{Annotations: map[string]string{"org.opencontainers.image.title": name}}); err == nil {
			t.Fatalf("accepted %s", name)
		}
	}
	if _, err := decodeManifest([]byte(`{"schemaVersion":2,"mediaType":"application/vnd.oci.image.index.v1+json"}`)); err == nil {
		t.Fatal("accepted index")
	}
}
func TestBlobMismatchRetriesAndCleansPartialFile(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { requests++; w.Write([]byte("corrupt")) }))
	defer server.Close()
	c := Client{HTTP: server.Client()}
	r := Reference{RegistryURL: server.URL, Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}
	d := Descriptor{Digest: digestOf([]byte("expected")), Size: 8}
	path := filepath.Join(t.TempDir(), "firmware")
	if err := retry(context.Background(), func() error { return c.downloadFile(context.Background(), r, d, path) }); err == nil {
		t.Fatal("accepted corruption")
	}
	if requests != 3 {
		t.Fatalf("requests=%d", requests)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("partial file remains: %v", err)
	}
}

func TestManifestMismatchAndMissingContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "manifests") {
			w.Header().Set("Docker-Content-Digest", "sha256:"+strings.Repeat("b", 64))
			w.Write([]byte(`{}`))
			return
		}
		w.WriteHeader(404)
	}))
	defer server.Close()
	c := Client{HTTP: server.Client()}
	if _, _, err := c.Resolve(context.Background(), server.URL, "firmware", "sha256:"+strings.Repeat("a", 64)); err == nil {
		t.Fatal("accepted mismatched manifest")
	}
	r := Reference{RegistryURL: server.URL, Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}
	m := Manifest{Config: Descriptor{Digest: r.ManifestDigest, Size: 1}}
	if err := c.EnsureAvailable(context.Background(), r, m); err == nil {
		t.Fatal("accepted missing blob")
	}
	if err := c.copyBlob(context.Background(), r, Descriptor{Digest: "invalid", Size: -1}, io.Discard); err == nil {
		t.Fatal("accepted invalid descriptor")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := retry(ctx, func() error { t.Fatal("called canceled operation"); return nil }); err == nil {
		t.Fatal("ignored cancellation")
	}
}
