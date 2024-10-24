/*******************************************************************************
* Copyright (c) Siemens AG 2022 ALL RIGHTS RESERVED.
*******************************************************************************/

package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchStringInSliceSmoketest(t *testing.T) {
	testCases := []struct {
		name   string
		list   []string
		search string
		found  bool
	}{
		{
			name:   "ElementNotFoundInsideList",
			list:   []string{"element1", "element2"},
			search: "not_found",
			found:  false,
		},
		{
			name:   "EmptySearchElement",
			list:   []string{"element1", "element2"},
			search: "",
			found:  false,
		},
		{
			name:   "ElementFound",
			list:   []string{"element1", "element2"},
			search: "element1",
			found:  true,
		},
		{
			name:   "SecondElementFound",
			list:   []string{"element1", "element2"},
			search: "element2",
			found:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, SearchStringInSlice(testCase.list, testCase.search), testCase.found)
		})
	}
}

func TestSearchSliceInSliceSmoketest(t *testing.T) {
	testCases := []struct {
		name   string
		list   []string
		search []string
		found  bool
	}{
		{
			name:   "NoElementFound",
			list:   []string{"element1", "element2"},
			search: []string{"not_found_1", "not_found_2"},
			found:  false,
		},
		{
			name:   "OneElementNotFound_FirstElement",
			list:   []string{"element1", "element2"},
			search: []string{"element1", "not_found"},
			found:  true,
		},
		{
			name:   "OneElementNotFound_SecondElement",
			list:   []string{"element1", "element2"},
			search: []string{"element2", "not_found"},
			found:  true,
		},
		{
			name:   "EmptySearchStringReturnNoElement",
			list:   []string{"element1", "element2"},
			search: []string{""},
			found:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, SearchSliceInSlice(testCase.list, testCase.search), testCase.found)
		})
	}
}
