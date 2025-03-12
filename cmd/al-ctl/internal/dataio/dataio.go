/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package dataio

import (
	"io"
	"os"
)

func ReadInput(filename string) ([]byte, error) {
	var fileReader = os.Stdin
	if filename != "" {
		var err error
		fileReader, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer fileReader.Close()
	}

	buffer, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func WriteOutput(filename string, data []byte) error {
	var fileWriter = os.Stdout
	if filename != "" {
		var err error
		fileWriter, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer fileWriter.Close()
	}

	_, err := fileWriter.Write(data)
	if err != nil {
		return err
	}
	return nil
}
