// +build gofuzz

package json

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"

	optimizedParser "github.com/mna/pigeon/examples/json/optimized"
	optimizedGrammar "github.com/mna/pigeon/examples/json/optimized-grammar"
)

var reUTF16Escape = regexp.MustCompile(`\\u[dD][89a-f-AF][0-9a-fA-F]{2}`)

// FuzzUnoptimized tests the unoptimized JSON parser.
func FuzzUnoptimized(fuzz []byte) int {
	_, err := Parse("fuzz", fuzz)
	if err != nil {
		return 0
	}
	return 1
}

// FuzzOptimizedParser tests the parser-optimized JSON parser.
func FuzzOptimizedParser(fuzz []byte) int {
	_, err := optimizedParser.Parse("fuzz", fuzz)
	if err != nil {
		return 0
	}
	return 1
}

// FuzzOptimizedGrammar tests the grammar-optimized JSON parser.
func FuzzOptimizedGrammar(fuzz []byte) int {
	_, err := optimizedGrammar.Parse("fuzz", fuzz)
	if err != nil {
		return 0
	}
	return 1
}

// FuzzInternalConsistency tests optimization variants against each other.
func FuzzInternalConsistency(fuzz []byte) int {
	unoptimizedResult, unoptimizedErr := Parse("fuzz", fuzz)
	opParserResult, opParserErr := optimizedParser.Parse("fuzz", fuzz)
	opGrammarResult, opGrammarErr := optimizedGrammar.Parse("fuzz", fuzz)
	if !reflect.DeepEqual(unoptimizedResult, opParserResult) {
		panic("bad parser-optimized result")
	}
	if !reflect.DeepEqual(unoptimizedErr, opParserErr) {
		panic("bad parser-optimized error")
	}
	if !reflect.DeepEqual(unoptimizedResult, opGrammarResult) {
		panic("bad grammar-optimized result")
	}
	if !reflect.DeepEqual(unoptimizedErr, opGrammarErr) {
		panic("bad grammar-optimized error")
	}
	if unoptimizedErr != nil {
		return 0
	}
	return 1
}

// FuzzExternalConsistency tests against the standard library parser.
func FuzzExternalConsistency(fuzz []byte) int {
	if reUTF16Escape.Match(fuzz) {
		return -1
	}
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
