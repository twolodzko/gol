package objects

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPush(t *testing.T) {
	var testCases = []struct {
		input    interface{}
		expected List
	}{
		{Int{42}, NewList(Int{42})},
		{String{"abc"}, NewList(Int{42}, String{"abc"})},
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
	expected := `(foo "bar" 42)`
	result := NewList(Symbol{"foo"}, String{"bar"}, Int{42}).String()

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}
