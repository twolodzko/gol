package parser

// import (
// 	"bufio"
// 	"reflect"
// 	"strings"
// 	"testing"
// )

// func TestReadList(t *testing.T) {
// 	var testCases = []struct {
// 		input    string
// 		expected []string
// 	}{
// 		{"(word)", []string{"word"}},
// 		{"(first second)", []string{"first", "second"}},
// 	}

// 	for _, test := range testCases {
// 		reader := bufio.NewReader(strings.NewReader(test.input))
// 		result, err := Parse(reader)

// 		if err != nil {
// 			t.Errorf("Unexpected error: %s", err)
// 		}
// 		if len(result) != len(test.expected) {
// 			t.Errorf("Lengths differ: %d vs %d", len(test.expected), len(result))
// 		}
// 		if !reflect.DeepEqual(result, test.expected) {
// 			t.Errorf("Experced: '%s' , got: '%s'", test.expected, result)
// 		}
// 	}
// }
