/*******************************************************************************

* Copyright (c) Siemens AG 2022 ALL RIGHTS RESERVED.

*******************************************************************************/

package shared

// SearchStringInSlice checks if a string is exactly inside an slice of strings
func SearchStringInSlice(list []string, searchString string) bool {
	for _, b := range list {
		if b == searchString {
			return true
		}
	}
	return false
}

// SearchSliceInSlice checks if one of the strings inside the slice is the slice of strings
func SearchSliceInSlice(list []string, searchStrings []string) bool {
	for _, value := range list {
		if SearchStringInSlice(searchStrings, value) {
			return true
		}
	}
	return false
}
