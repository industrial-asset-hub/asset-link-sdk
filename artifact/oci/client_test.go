// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package oci

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const payloadTitleAnnotation = "org.opencontainers.image.title"

func TestResolveEnsuresBlobsAndDownloadsPayload(t *testing.T) {
	payload := []byte("firmware bytes")
	config := []byte("{}")
	manifest := Manifest{
		SchemaVersion: 2,
		MediaType:     ManifestMediaType,
		Config:        descriptorFor(config),
		Layers: []Descriptor{
			{
				Digest: digestOf(payload),
				Size:   int64(len(payload)),
				Annotations: map[string]string{
					payloadTitleAnnotation: "firmware.upd",
				},
			},
		},
	}
	manifestData, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}

	blobs := map[string][]byte{
		digestOf(config):  config,
		digestOf(payload): payload,
	}
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/manifests/"):
			writeResponse(t, w, manifestData)
		case strings.Contains(r.URL.Path, "/blobs/"):
			writeResponse(t, w, blobs[filepath.Base(r.URL.Path)])
		default:
			http.NotFound(w, r)
		}
	})
	client := Client{HTTP: server.Client()}
	ctx := context.Background()

	reference, resolvedManifest, err := client.Resolve(ctx, server.URL, "firmware/s7-1500", "latest")
	if err != nil {
		t.Fatalf("resolve manifest: %v", err)
	}
	if want := digestOf(manifestData); reference.ManifestDigest != want {
		t.Errorf("manifest digest = %q, want %q", reference.ManifestDigest, want)
	}

	if err := client.EnsureAvailable(ctx, reference, resolvedManifest); err != nil {
		t.Fatalf("ensure blobs are available: %v", err)
	}

	directory := t.TempDir()
	files, err := client.Download(ctx, reference, directory)
	if err != nil {
		t.Fatalf("download firmware: %v", err)
	}
	if want := []string{"firmware.upd"}; !slices.Equal(files, want) {
		t.Fatalf("downloaded files = %v, want %v", files, want)
	}

	gotPayload, err := os.ReadFile(filepath.Join(directory, "firmware.upd"))
	if err != nil {
		t.Fatalf("read downloaded firmware: %v", err)
	}
	if !bytes.Equal(gotPayload, payload) {
		t.Errorf("downloaded payload = %q, want %q", gotPayload, payload)
	}
}

func TestRepositoryNames(t *testing.T) {
	for _, repository := range []string{"firmware", "vendor/firmware.v1", "vendor/my_firmware", "vendor/my__firmware", "vendor/my--firmware", "vendor/my---firmware"} {
		t.Run(repository, func(t *testing.T) {
			reference := testReference("https://registry.example")
			reference.Repository = repository
			if err := reference.Validate(); err != nil {
				t.Fatalf("valid repository rejected: %v", err)
			}
		})
	}
	for _, repository := range []string{"", "/firmware", "firmware/", "vendor//firmware", "Firmware", "vendor/my___firmware", "vendor/my..firmware", "vendor/my_-firmware", "vendor/-firmware", "vendor/firmware-"} {
		t.Run(repository, func(t *testing.T) {
			reference := testReference("https://registry.example")
			reference.Repository = repository
			if err := reference.Validate(); err == nil {
				t.Fatal("invalid repository accepted")
			}
		})
	}
}

func TestResolveRepositorySeparators(t *testing.T) {
	manifest, err := json.Marshal(Manifest{
		SchemaVersion: 2,
		MediaType:     ManifestMediaType,
		Config:        descriptorFor([]byte("{}")),
		Layers:        []Descriptor{descriptorFor([]byte("firmware"))},
	})
	if err != nil {
		t.Fatal(err)
	}
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		writeResponse(t, w, manifest)
	})
	client := Client{HTTP: server.Client()}
	for _, repository := range []string{"vendor/my__firmware", "vendor/my--firmware"} {
		t.Run(repository, func(t *testing.T) {
			if _, _, err := client.Resolve(context.Background(), server.URL, repository, "latest"); err != nil {
				t.Fatalf("resolve valid repository: %v", err)
			}
		})
	}
}

func TestRejectUnsafeRegistryURLs(t *testing.T) {
	registryURLs := []string{
		"file:///tmp",
		"http://user:password@localhost",
		"http://localhost/path",
		"http://localhost?x=1",
	}
	for _, registryURL := range registryURLs {
		t.Run(registryURL, func(t *testing.T) {
			reference := testReference(registryURL)
			if err := reference.Validate(); err == nil {
				t.Fatalf("Validate() accepted unsafe registry URL %q", registryURL)
			}
		})
	}
}

func TestRejectUnsafePayloadNames(t *testing.T) {
	payloadNames := []string{"../escape", "sub/file", `sub\file`, "..", "prepared.json"}
	for _, name := range payloadNames {
		t.Run(name, func(t *testing.T) {
			descriptor := Descriptor{
				Annotations: map[string]string{payloadTitleAnnotation: name},
			}
			if _, err := payloadName(descriptor); err == nil {
				t.Fatalf("payloadName() accepted unsafe name %q", name)
			}
		})
	}
}

func TestRejectManifestIndex(t *testing.T) {
	index := []byte(`{"schemaVersion":2,"mediaType":"application/vnd.oci.image.index.v1+json"}`)
	if _, err := decodeManifest(index); err == nil {
		t.Fatal("decodeManifest() accepted an OCI image index")
	}
}

func TestBlobMismatchRetriesAndCleansPartialFile(t *testing.T) {
	requests := 0
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		requests++
		writeResponse(t, w, []byte("corrupt"))
	})
	client := Client{HTTP: server.Client()}
	reference := testReference(server.URL)
	descriptor := descriptorFor([]byte("expected"))
	path := filepath.Join(t.TempDir(), "firmware")

	err := retry(context.Background(), func() error {
		return client.downloadFile(context.Background(), reference, descriptor, path)
	})
	if err == nil {
		t.Fatal("download succeeded despite the blob digest mismatch")
	}
	if want := 3; requests != want {
		t.Errorf("blob requests = %d, want %d", requests, want)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("partial download should be removed; stat error = %v", err)
	}
}

func TestRejectManifestDigestMismatch(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/manifests/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Docker-Content-Digest", digestOf([]byte("different manifest")))
		writeResponse(t, w, []byte(`{}`))
	})
	client := Client{HTTP: server.Client()}
	manifestDigest := "sha256:" + strings.Repeat("a", 64)

	_, _, err := client.Resolve(context.Background(), server.URL, "firmware", manifestDigest)
	if err == nil {
		t.Fatal("Resolve() accepted a mismatched manifest digest")
	}
	if !strings.Contains(err.Error(), "digest mismatch") {
		t.Fatalf("Resolve() error = %v, want a digest mismatch", err)
	}
}

func TestEnsureAvailableRejectsMissingBlob(t *testing.T) {
	server := newTestServer(t, http.NotFound)
	client := Client{HTTP: server.Client()}
	reference := testReference(server.URL)
	manifest := Manifest{Config: descriptorFor([]byte("config"))}

	err := client.EnsureAvailable(context.Background(), reference, manifest)
	if err == nil {
		t.Fatal("EnsureAvailable() accepted a missing blob")
	}
}

func TestCopyBlobRejectsInvalidDescriptor(t *testing.T) {
	client := Client{}
	reference := testReference("https://registry.example")
	descriptor := Descriptor{Digest: "invalid", Size: -1}

	if err := client.copyBlob(context.Background(), reference, descriptor, io.Discard); err == nil {
		t.Fatal("copyBlob() accepted an invalid descriptor")
	}
}

func TestRetryStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := retry(ctx, func() error {
		t.Fatal("retry operation called with a canceled context")
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("retry() error = %v, want %v", err, context.Canceled)
	}
}

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func writeResponse(t *testing.T, w http.ResponseWriter, body []byte) {
	t.Helper()
	if _, err := w.Write(body); err != nil {
		t.Errorf("write test response: %v", err)
	}
}

func testReference(registryURL string) Reference {
	return Reference{
		RegistryURL:    registryURL,
		Repository:     "firmware",
		ManifestDigest: "sha256:" + strings.Repeat("a", 64),
	}
}

func descriptorFor(data []byte) Descriptor {
	return Descriptor{Digest: digestOf(data), Size: int64(len(data))}
}
