/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package apimock

import (
	"fmt"
)

type testFuncs func(string, string) bool

func RunApiMockTests(address, discoveryFile string) {
	// Add all the test functions here
	allTests := []testFuncs{
		TestStartDiscovery,
		TestGetFilterTypes,
		TestGetFilterOptions,
	}
	testPassed := 0
	for _, test := range allTests {
		if !test(address, discoveryFile) {
			fmt.Println("Test failed")
			continue
		}
		fmt.Println("Test passed")
		testPassed++
	}
	fmt.Printf("Total tests passed: %d/%d, failed: %d\n", testPassed, len(allTests), len(allTests)-testPassed)
}
