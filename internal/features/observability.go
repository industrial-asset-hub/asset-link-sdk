/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package features

type observabilityFeatures struct {
	HttpObservabilityServer bool
}

var config = ObservabilityFeaturesNew()

func ObservabilityFeatures() *observabilityFeatures {
	return config
}

func ObservabilityFeaturesNew() *observabilityFeatures {
	return &observabilityFeatures{
		HttpObservabilityServer: false,
	}
}

func (o *observabilityFeatures) Get() *observabilityFeatures {
	return o
}
