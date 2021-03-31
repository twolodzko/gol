package evaluator

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/parser"

	. "github.com/twolodzko/goal/types"
)

func TestEval(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Any
	}{
		// objects
		{"()", List{}},
		{"2", Int(2)},
		{"3.14", Float(3.14)},
		{`"Hello World!"`, String("Hello World!")},
		// {"true", Bool(true)},
		// functions
		{`(str 3.14)`, String("3.14")},
		{`(int "3.14")`, Int(3)},
		{`(float "1")`, Float(1)},
		{`(list "Hello World!" 42 3.14)`, List{String("Hello World!"), Int(42), Float(3.14)}},
		{`(quote 3.14)`, Float(3.14)},
		{`(quote foo)`, Symbol("foo")},
		{`(quote (foo bar))`, List{Symbol("foo"), Symbol("bar")}},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := Eval(expr[0])

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

// func TestBooleans(t *testing.T) {
// 	var testCases = []struct {
// 		input    string
// 		expected Bool
// 	}{
// 		// booleans: everything is true, unless it's false
// 		{`(true? "1")`, true},
// 		{`(true? 0)`, true},
// 		{`(true? 3.1415)`, true},
// 		{`(true? true)`, true},
// 		{`(true? ())`, true},
// 		{`(true? false)`, false},
// 		{`(not true)`, false},
// 		{`(not false)`, true},
// 		{`(not ())`, false},
// 	}

// 	for _, tt := range testCases {
// 		expr, err := parser.Parse(strings.NewReader(tt.input))

// 		if err != nil {
// 			t.Errorf("unexpected error: %s", err)
// 		}

// 		result, err := Eval(expr[0])

// 		if err != nil {
// 			t.Errorf("unexpected error: %s", err)
// 		}
// 		if !cmp.Equal(result, tt.expected) {
// 			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
// 		}
// 	}
// }
