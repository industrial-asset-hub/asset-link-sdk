/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	_ "errors"
	"fmt"
	"regexp"
)

var (
	ErrValidation = &ValidationError{}
	ErrEmpty      = &EmptyError{}
)

// ValidationError represents a validation failure with optional details.
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
	Details interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// EmptyError represents an error for empty required values.
type EmptyError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e *EmptyError) Error() string {
	return fmt.Sprintf("%s (field: '%s', value: %v)", e.Message, e.Field, e.Value)
}

// PermissibleValuesError represents an error for values not in the allowed set.
type PermissibleValuesError struct {
	Field   string
	Message string
	Value   interface{}
	Allowed []interface{}
}

func (e *PermissibleValuesError) Error() string {
	return fmt.Sprintf("field '%s' has value '%v' which is not in permissible values %v", e.Field, e.Value, e.Allowed)
}

// ValidateField checks if a value is non-empty and matches a pattern, returning the appropriate error or nil.
func ValidateField(value string, fieldName, emptyMsg, pattern, patternMsg string) error {
	if !isNonEmptyValues(value) {
		return &EmptyError{
			Field:   fieldName,
			Message: emptyMsg,
			Value:   value,
		}
	}
	if pattern != "" && !ValidateByPattern(value, pattern) {
		return &ValidationError{
			Field:   fieldName,
			Message: patternMsg,
			Value:   value,
			Details: value,
		}
	}
	return nil
}

func ValidateByPattern(value, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(value)
}
