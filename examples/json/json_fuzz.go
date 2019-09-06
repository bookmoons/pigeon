// +build gofuzz

package json

// Fuzz tests the JSON parser.
func Fuzz(fuzz []byte) int {
	_, err := Parse("fuzz", fuzz)
	if err != nil {
		return 0
	}
	return 1
}
