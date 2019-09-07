// +build gofuzz

package json

import (
	"encoding/json"
	"reflect"
)

// Fuzz tests the JSON parser.
func Fuzz(fuzz []byte) int {
	var expected interface{}
	expectedErr := json.Unmarshal(fuzz, &expected)
	result, err := Parse("fuzz", fuzz)
	if err != nil {
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
