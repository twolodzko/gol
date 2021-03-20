package list

import (
	"testing"
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
		same, err := AreSame(tt.x, tt.y)

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
		{42, NewList(42)},
		{"abc", NewList(42, "abc")},
	}

	list := List{}
	for _, tt := range testCases {
		list.Push(tt.input)
		if same, err := AreSame(list, tt.expected); !same {
			t.Errorf("expected: %v, got: %v", tt.expected, list)

			if err != nil {
				t.Errorf("%v", err)
			}
		}
	}
}
