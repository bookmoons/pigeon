// +build gofuzz

package json

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Fuzz tests the JSON parser.
func Fuzz(fuzz []byte) int {
	var expected interface{}
	expectedErr := json.Unmarshal(fuzz, &expected)
	result, err := Parse("fuzz", fuzz)
	if err != nil {
		if strings.Contains(err.Error(), "invalid encoding") {
			return -1
		}
		if expectedErr == nil {
			panic("incorrect failure")
		}
		return 0
	} else if expectedErr != nil {
		panic("incorrect success")
	}
	if !reflect.DeepEqual(result, expected) {
		panic("unexpected result")
	}
	return 1
}
