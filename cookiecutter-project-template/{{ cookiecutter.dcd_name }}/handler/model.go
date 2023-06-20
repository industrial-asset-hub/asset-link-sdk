/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 */

package handler

// Default Device Information structure
type DeviceInfo struct {
  Vendor       string `json:"vendor"`
  SerialNumber string `json:"serial_number"`
  DeviceFamily string `json:"device_family"`
}
