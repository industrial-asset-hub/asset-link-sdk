// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

// Package oci retrieves firmware artifacts over HTTP, not gRPC.
package oci

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const ManifestMediaType = "application/vnd.oci.image.manifest.v1+json"

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
var repositoryPattern = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*(?:/[a-z0-9]+(?:[._-][a-z0-9]+)*)*$`)
var tagPattern = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}$`)

type Reference struct {
	RegistryURL    string `json:"registryUrl"`
	Repository     string `json:"repository"`
	ManifestDigest string `json:"manifestDigest"`
}
type Descriptor struct {
	Digest      string            `json:"digest"`
	Size        int64             `json:"size"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
type Manifest struct {
	SchemaVersion int          `json:"schemaVersion"`
	MediaType     string       `json:"mediaType"`
	Config        Descriptor   `json:"config"`
	Layers        []Descriptor `json:"layers"`
}
type Client struct{ HTTP *http.Client }

func (r Reference) Validate() error {
	if _, err := registryBase(r.RegistryURL); err != nil {
		return err
	}
	if !repositoryPattern.MatchString(r.Repository) || !digestPattern.MatchString(r.ManifestDigest) {
		return fmt.Errorf("OCI reference requires a repository and sha256 manifest digest")
	}
	return nil
}
func registryBase(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
		return "", fmt.Errorf("invalid registry URL")
	}
	if (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("registry URL must not contain path, query or fragment")
	}
	return strings.TrimRight(raw, "/"), nil
}

// Resolve triggers Zot synchronization and pins a tag to its manifest digest.
func (c Client) Resolve(ctx context.Context, registry, repository, selection string) (Reference, Manifest, error) {
	r := Reference{RegistryURL: registry, Repository: repository}
	if !repositoryPattern.MatchString(repository) || (!digestPattern.MatchString(selection) && !tagPattern.MatchString(selection)) {
		return r, Manifest{}, fmt.Errorf("invalid OCI repository or selection")
	}
	base, err := registryBase(registry)
	if err != nil {
		return r, Manifest{}, err
	}
	response, err := c.get(ctx, base+"/v2/"+repository+"/manifests/"+selection)
	if err != nil {
		return r, Manifest{}, err
	}
	defer response.Body.Close()
	data, err := readManifest(response, selection)
	if err != nil {
		return r, Manifest{}, err
	}
	r.ManifestDigest = digestOf(data)
	m, err := decodeManifest(data)
	return r, m, err
}
func readManifest(response *http.Response, selection string) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 4<<20 {
		return nil, fmt.Errorf("OCI manifest exceeds 4 MiB")
	}
	digest := digestOf(data)
	if digestPattern.MatchString(selection) && selection != digest {
		return nil, fmt.Errorf("manifest digest mismatch")
	}
	if header := response.Header.Get("Docker-Content-Digest"); header != "" && header != digest {
		return nil, fmt.Errorf("registry manifest digest mismatch")
	}
	return data, nil
}
func decodeManifest(data []byte) (Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return m, err
	}
	if m.SchemaVersion != 2 || m.MediaType != ManifestMediaType || len(m.Layers) == 0 {
		return m, fmt.Errorf("firmware requires a single OCI image manifest with payload layers (indexes unsupported)")
	}
	for _, d := range append([]Descriptor{m.Config}, m.Layers...) {
		if !digestPattern.MatchString(d.Digest) || d.Size < 0 || d.Size > 1<<40 {
			return m, fmt.Errorf("invalid OCI content descriptor")
		}
	}
	return m, nil
}
func (c Client) get(ctx context.Context, address string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", ManifestMediaType)
	request.Header.Set("Accept-Encoding", "identity")
	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("OCI GET returned HTTP %d", response.StatusCode)
	}
	return response, nil
}
func digestOf(data []byte) string {
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:])
}

// EnsureAvailable checks every blob using bounded memory, not just the manifest.
func (c Client) EnsureAvailable(ctx context.Context, r Reference, m Manifest) error {
	if err := r.Validate(); err != nil {
		return err
	}
	for _, d := range append([]Descriptor{m.Config}, m.Layers...) {
		if err := retry(ctx, func() error { return c.copyBlob(ctx, r, d, io.Discard) }); err != nil {
			return err
		}
	}
	return nil
}
func (c Client) copyBlob(ctx context.Context, r Reference, d Descriptor, writer io.Writer) error {
	if !digestPattern.MatchString(d.Digest) || d.Size < 0 || d.Size > 1<<40 {
		return fmt.Errorf("invalid blob descriptor")
	}
	response, err := c.get(ctx, strings.TrimRight(r.RegistryURL, "/")+"/v2/"+r.Repository+"/blobs/"+d.Digest)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	hash := sha256.New()
	n, err := io.Copy(io.MultiWriter(writer, hash), io.LimitReader(response.Body, d.Size+1))
	if err != nil {
		return err
	}
	if n != d.Size || "sha256:"+hex.EncodeToString(hash.Sum(nil)) != d.Digest {
		return fmt.Errorf("blob size or digest mismatch")
	}
	return nil
}
func retry(ctx context.Context, operation func() error) error {
	var err error
	for range 3 {
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = operation(); err == nil {
			return nil
		}
	}
	return err
}

// Download stages verified payload files without extracting archives. Use a
// private temporary directory and publish that directory only on success.
func (c Client) Download(ctx context.Context, r Reference, directory string) ([]string, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	_, m, err := c.Resolve(ctx, r.RegistryURL, r.Repository, r.ManifestDigest)
	if err != nil {
		return nil, err
	}
	var files []string
	seen := map[string]bool{}
	for _, d := range m.Layers {
		name, err := payloadName(d)
		if err != nil {
			return nil, err
		}
		if seen[name] {
			return nil, fmt.Errorf("duplicate firmware payload name %q", name)
		}
		seen[name] = true
		path := filepath.Join(directory, name)
		if err := retry(ctx, func() error { return c.downloadFile(ctx, r, d, path) }); err != nil {
			return nil, err
		}
		files = append(files, name)
	}
	return files, nil
}
func payloadName(d Descriptor) (string, error) {
	name := d.Annotations["org.opencontainers.image.title"]
	if name == "" {
		return strings.TrimPrefix(d.Digest, "sha256:") + ".bin", nil
	}
	if name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00") || name == "prepared.json" {
		return "", fmt.Errorf("unsafe firmware payload name")
	}
	return name, nil
}
func (c Client) downloadFile(ctx context.Context, r Reference, d Descriptor, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	err = c.copyBlob(ctx, r, d, file)
	closeErr := file.Close()
	if err != nil {
		os.Remove(path)
		return err
	}
	return closeErr
}
