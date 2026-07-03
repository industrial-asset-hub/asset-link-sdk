/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
)

func TestGetDeviceInfoPropertiesEmptyRequestFilePath(t *testing.T) {
	responses, err := GetDeviceInfoProperties("localhost:50051", "")
	if err == nil {
		t.Fatal("expected error for empty request file path")
	}
	if responses != nil {
		t.Fatal("expected nil responses for empty request file path")
	}
}

func TestCreatePropertyValuesRequestFromInputFileSuccess(t *testing.T) {
	tempDir := t.TempDir()
	requestFile := filepath.Join(tempDir, "request.json")

	if err := os.WriteFile(requestFile, []byte(`{"keys":["name","serialNumber"]}`), 0o600); err != nil {
		t.Fatalf("failed to write request file: %v", err)
	}

	request, err := createPropertyValuesRequestFromInputFile(requestFile)
	if err != nil {
		t.Fatalf("unexpected error creating request from file: %v", err)
	}
	if request == nil {
		t.Fatal("expected non-nil request")
	}
	if got := len(request.GetKeys()); got != 2 {
		t.Fatalf("unexpected key count: got %d want 2", got)
	}
}

func TestCreatePropertyValuesRequestFromInputFileMissingFile(t *testing.T) {
	request, err := createPropertyValuesRequestFromInputFile("missing-request.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if request != nil {
		t.Fatal("expected nil request for missing file")
	}
}

func TestCreatePropertyValuesRequestFromInputFileInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	requestFile := filepath.Join(tempDir, "invalid-request.json")

	if err := os.WriteFile(requestFile, []byte(`{"keys":[`), 0o600); err != nil {
		t.Fatalf("failed to write invalid request file: %v", err)
	}

	request, err := createPropertyValuesRequestFromInputFile(requestFile)
	if err == nil {
		t.Fatal("expected error for invalid request json")
	}
	if request != nil {
		t.Fatal("expected nil request for invalid request json")
	}
}

func TestWritePropertyResponsesFileNilResponse(t *testing.T) {
	err := WritePropertyResponsesFile(filepath.Join(t.TempDir(), "response.json"), nil)
	if err == nil {
		t.Fatal("expected error for nil property response")
	}
}

func TestWritePropertyResponsesFileSuccess(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "response.json")

	responses := &generatedDeviceInfo.GetPropertyValuesResponse{}
	if err := WritePropertyResponsesFile(outputFile, responses); err != nil {
		t.Fatalf("unexpected error writing property responses: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty output file")
	}

	var decoded map[string]any
	if err := json.Unmarshal(content, &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}
