package parser

import (
	"reflect"
	"testing"
)

func TestPush(t *testing.T) {
	var testCases = []struct {
		input    interface{}
		expected List
	}{
		{42, NewList(42)},
		//{"abc", NewList(42, "abc")},
	}

	list := List{}
	for _, tt := range testCases {
		list.Push(tt.input)
		if !reflect.DeepEqual(list, tt.expected) {
			t.Errorf("experted: '%s' , got: '%s'", tt.expected, list)
		}
	}
}

func TestParse(t *testing.T) {
	var testCases = []struct {
		input    string
		expected List
	}{
		{"", List{}},
		{"42", NewList(42)},
	}

	for _, tt := range testCases {
		result, err := Parse(tt.input)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("experted: '%s' , got: '%s'", tt.expected, result)
		}
	}
}
