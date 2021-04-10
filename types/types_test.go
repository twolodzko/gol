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

	if input.String() != expected {
		t.Errorf("expected %s, got: %s", expected, input.String())
	}

	str := String("abc\n\t\"")
	if str.Raw() != "abc\n\t\"" {
		t.Errorf("expected %s, got: %s", "abc\n\t\"", str.Raw())
	}
	if str.Quote() != `abc\n\t\"` {
		t.Errorf("expected %s, got: %s", `abc\n\t\"`, str.Quote())
	}
	if str.String() != "\"abc\n\t\"\"" {
		t.Errorf("expected %s, got: %s", "\"abc\n\t\"\"", str.String())
	}

	unquoted, err := str.Quote().Unquote()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if unquoted != str {
		t.Errorf("expected %s, got: %s", str, unquoted)
	}
}
