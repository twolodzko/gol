package objects

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPush(t *testing.T) {
	var testCases = []struct {
		input    Object
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

func TestListHead(t *testing.T) {
	var testCases = []struct {
		input    List
		expected Object
	}{
		{NewList(Int{1}), Int{1}},
		{NewList(Int{1}, Int{2}, Int{3}), Int{1}},
	}

	for _, tt := range testCases {
		result := tt.input.Head()

		if !cmp.Equal(tt.expected, result) {
			t.Errorf("expected: %v, got: %v", tt.expected, result)
		}
	}
}

func TestListTail(t *testing.T) {
	var testCases = []struct {
		input    List
		expected Object
	}{
		{NewList(Int{1}), List{}},
		{NewList(Int{1}, Int{2}, Int{3}), NewList(Int{2}, Int{3})},
	}

	for _, tt := range testCases {
		result := tt.input.Tail()

		if !cmp.Equal(tt.expected, result) {
			t.Errorf("expected: %v, got: %v", tt.expected, result)
		}
	}
}
