package old_parser

import (
	"reflect"
	"testing"
)

func TestTrivial(t *testing.T) {
	input := "(def pi 3.14)"
	expected := []string{"def", "pi", "3.14"}
	result, err := Parse(input)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got: %s, expected: %s", result, expected)
	}
}

func TestMissingBracket(t *testing.T) {
	input := "(def pi 3.14"
	_, err := Parse(input)

	if err == nil {
		t.Errorf("Expected to see an error")
	}
}
