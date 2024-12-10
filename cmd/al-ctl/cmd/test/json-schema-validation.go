/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJsonSchema(schemaPath string, dataPath string) error {
	// Get absolute path for schema
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return fmt.Errorf("could not get absolute path for schema: %w", err)
	}
	schemaURL := "file://" + strings.ReplaceAll(absSchemaPath, "\\", "/")

	// Get absolute path for data
	absDataPath, err := filepath.Abs(dataPath)
	if err != nil {
		return fmt.Errorf("could not get absolute path for data: %w", err)
	}
	dataURL := "file://" + strings.ReplaceAll(absDataPath, "\\", "/")

	// Validate
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)
	dataLoader := gojsonschema.NewReferenceLoader(dataURL)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return fmt.Errorf("error validating JSON: %w", err)
	}

	if !result.Valid() {
		errorMsg := "The JSON data is invalid. Errors:"
		for _, err := range result.Errors() {
			errorMsg += fmt.Sprintf("\n- %s", err)
		}
		return fmt.Errorf("%s", errorMsg)
	}

	fmt.Println("The JSON data is valid")
	return nil
}
