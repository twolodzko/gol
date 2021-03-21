package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAreSame(t *testing.T) {
	var testCases = []struct {
		x    List
		y    List
		same bool
	}{
		{List{}, List{}, true},
		{List{[]interface{}{42}}, List{[]interface{}{42}}, true},
		{List{[]interface{}{1, 2}}, List{[]interface{}{1, 2}}, true},
		{List{[]interface{}{1, "abc"}}, List{[]interface{}{1, "abc"}}, true},
		{List{[]interface{}{1, 2}}, List{[]interface{}{1}}, false},
		{List{[]interface{}{1, 2}}, List{[]interface{}{1, "2"}}, false},
	}

	for _, tt := range testCases {
		same, err := areSame(tt.x, tt.y)

		if same != tt.same {
			t.Errorf("failed comparison for lists %v and %v", tt.x, tt.y)
		}
		if tt.same && err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

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
