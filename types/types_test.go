package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestList(t *testing.T) {
	var testCases = []struct {
		input List
		head  Any
		tail  List
	}{
		{List{}, nil, List{}},
		{List{1}, 1, List{}},
		{List{1, 2, 3}, 1, List{2, 3}},
		{List{List{}}, List{}, List{}},
		{List{List{}, List{}}, List{}, List{List{}}},
	}

	for _, tt := range testCases {
		head := tt.input.Head()
		if !cmp.Equal(head, tt.head) {
			t.Errorf("expected %v (%T), got: %v (%T)", tt.head, tt.head, head, head)
		}

		tail := tt.input.Tail()
		if !cmp.Equal(tail, tt.tail) {
			t.Errorf("expected %v (%T), got: %v (%T)", tt.tail, tt.tail, tail, tail)
		}
	}
}

func TestString(t *testing.T) {
	input := List{true, 42, 3.14, Symbol("foo"), String("Hello World!")}
	expected := `(true 42 3.14 foo "Hello World!")`
	result := input.String()

	if result != expected {
		t.Errorf("expected %s, got: %s", expected, result)
	}
}
