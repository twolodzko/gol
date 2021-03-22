package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPush(t *testing.T) {
	var testCases = []struct {
		input    interface{}
		expected List
	}{
		{42, newList(42)},
		{"abc", newList(42, "abc")},
	}

	list := List{}
	for _, tt := range testCases {
		list.Push(tt.input)

		if !cmp.Equal(list, tt.expected) {
			t.Errorf("expected: %v, got: %v", tt.expected, list)
		}
	}
}

func TestString(t *testing.T) {
	expected := "(foo \"bar\" 42)"
	result := newList(Symbol{"foo"}, String{"bar"}, 42).String()

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}
