package evaluator

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/parser"

	. "github.com/twolodzko/goal/types"
)

func TestEvalExpr(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Any
	}{
		// objects
		{`()`, List{}},
		{`2`, Int(2)},
		{`3.14`, Float(3.14)},
		{`"Hello World!"`, String("Hello World!")},
		{`true`, Bool(true)},
		// functions
		{`(str 3.14)`, String("3.14")},
		{`(int "3.14")`, Int(3)},
		{`(float "1")`, Float(1)},
		{`(list "Hello World!" 42 3.14)`, List{String("Hello World!"), Int(42), Float(3.14)}},
		{`(quote 3.14)`, Float(3.14)},
		{`(quote foo)`, Symbol("foo")},
		{`(quote (foo bar))`, List{Symbol("foo"), Symbol("bar")}},
		{`(size ())`, Int(0)},
		{`(size (list 1 2 3))`, Int(3)},
		{`(size "")`, Int(0)},
		{`(size "hello")`, Int(5)},
		{`(size (list 1 2 3) () (quote foo bar) "abcd")`, List{Int(3), Int(0), Int(2), Int(4)}},
		{`(head (list 1 2 3))`, Int(1)},
		{`(tail (list 1 2 3))`, List{Int(2), Int(3)}},
		{`(nil? nil)`, Bool(true)},
		{`(nil? ())`, Bool(false)},
		{`(nil? true)`, Bool(false)},
		{`(if true 1 2)`, Int(1)},
		{`(if false 1 2)`, Int(2)},
		{`(if (true? false) (int+ 2 2) (int- 2 2))`, Int(0)},
		{`(= 2 2)`, Bool(true)},
		{`(= 2 3)`, Bool(false)},
		{`(= 2 "2")`, Bool(false)},
		{`(= 2 (int+ 1 1))`, Bool(true)},
		// math
		{`(int+ 2 2)`, Int(4)},
		{`(int+ 2 2 2 2)`, Int(8)},
		{`(int- 3 2)`, Int(1)},
		{`(int* 2 3)`, Int(6)},
		{`(int/ 6 3)`, Int(2)},
		{`(float+ 2.1 4.15)`, Float(6.25)},
		{`(float- 2.1 4.0)`, Float(-1.9)},
		{`(float* 2.5 4.0)`, Float(10.0)},
		{`(float/ 10.2 5.1)`, Float(2.0)},
		{`(int+ 2 (int- 4 (int* 1 2)))`, Int(4)},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := EvalExpr(expr[0], BaseEnv)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}

func TestBooleans(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Bool
	}{
		// booleans: everything is true
		// the only false things are Bool(false) and nil
		{`(true? "1")`, true},
		{`(true? 0)`, true},
		{`(true? 3.1415)`, true},
		{`(true? true)`, true},
		{`(true? ())`, true},
		{`(true? false)`, false},
		{`(not true)`, false},
		{`(not false)`, true},
		{`(not ())`, false},
		{`(true? nil)`, false},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := EvalExpr(expr[0], BaseEnv)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func TestErrorFn(t *testing.T) {
	input := `(list 1 (error "ok!") 2)`
	expr, err := parser.Parse(strings.NewReader(input))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	result, err := EvalExpr(expr[0], BaseEnv)

	if err == nil {
		t.Errorf("expected error, got result: %v", result)
	}
}

func TestLet(t *testing.T) {
	expr, err := parser.Parse(strings.NewReader(`(let x 42)`))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if _, err := EvalExpr(expr[0], BaseEnv); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if _, err := BaseEnv.Get(Symbol("x")); err != nil {
		t.Errorf("variable x not set")
	}

	expr, err = parser.Parse(strings.NewReader("x"))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	result, err := EvalExpr(expr[0], BaseEnv)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != Int(42) {
		t.Errorf("unable to read the variable")
	}
}
