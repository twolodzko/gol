package evaluator

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/environment"
	"github.com/twolodzko/goal/parser"

	. "github.com/twolodzko/goal/types"
)

func TestEval(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Any
	}{
		{`nil`, nil},
		{`()`, List{}},
		{`2`, Int(2)},
		{`3.14`, Float(3.14)},
		{`"Hello World!"`, String("Hello World!")},
		{`true`, Bool(true)},
		{`(if true 1 2)`, Int(1)},
		{`(if false 1 2)`, Int(2)},
		{`(if (true? false) (error "this should not fail!") "ok")`, String("ok")},
		{`(quote (+ 1 2))`, List{Symbol("+"), Int(1), Int(2)}},
		{`(- 7 (* 2 (+ 1 2)) 1)`, Int(0)},
		{`(def b (+ 1 2))`, Int(3)},
		{`(let (c 2) c)`, Int(2)},
		{`(let (x 10) (+ 5 x))`, Int(15)},
		{`(let (x 5) (let (y 6) (+ x y)))`, Int(11)},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		env := InitBuildin()
		result, err := Eval(expr[0], env)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}

func TestCore(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Any
	}{
		// {`(str 3.14 42 "hello")`, List{String("3.14"), String("42"), String("hello")}},
		// {`(int "3.14" "10" 5.2 100)`, List{Int(3), Int(10), Int(5), Int(100)}},
		// {`(float 5.22 "1" "1e-5")`, List{Float(5.22), Float(1), Float(1e-5)}},
		{`(list "Hello World!" 42 3.14)`, List{String("Hello World!"), Int(42), Float(3.14)}},
		{`(quote 3.14)`, Float(3.14)},
		{`(quote foo)`, Symbol("foo")},
		{`(quote (foo bar))`, List{Symbol("foo"), Symbol("bar")}},
		{`(head (list 1 2 3))`, Int(1)},
		{`(tail (list 1 2 3))`, List{Int(2), Int(3)}},
		{`(eq? 2 2)`, Bool(true)},
		{`(eq? 2 3)`, Bool(false)},
		{`(eq? 2 "2")`, Bool(false)},
		{`(eq? (list 1 2 3) (list 1 2 3))`, Bool(true)},
		{`(eq? (list 1 2 3) (list 1 "2" 3))`, Bool(false)},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		env := InitBuildin()
		result, err := Eval(expr[0], env)

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
		{`(and true true false)`, false},
		{`(and true 1 ())`, true},
		{`(or false true false)`, true},
		{`(or false nil)`, false},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		env := InitBuildin()
		result, err := Eval(expr[0], env)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func TestMath(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Any
	}{
		{`(+ 2 2.0)`, Float(4.0)},
		{`(- 3 2)`, Int(1)},
		{`(* 2 3)`, Int(6)},
		{`(/ 6 3)`, Float(2)},
		{`(+ 2.1 4.15)`, Float(6.25)},
		{`(- 2.1 4.0)`, Float(-1.9)},
		{`(* 2.5 4.0)`, Float(10.0)},
		{`(/ 10.2 5.1)`, Float(2.0)},
		{`(+ 2 (- 4 (* 1 2)))`, Int(4)},
		{`(// 100 2 5 2)`, Int(5)},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		env := InitBuildin()
		result, err := Eval(expr[0], env)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}

func TestErrorFn(t *testing.T) {
	input := `(list 1 (error "ok!") 2)`
	expr, err := parser.Parse(strings.NewReader(input))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	env := environment.NewEnv()
	result, err := Eval(expr[0], env)

	if err == nil {
		t.Errorf("expected error, got result: %v", result)
	}
}

func TestDef(t *testing.T) {
	env := InitBuildin()

	expr, err := parser.Parse(strings.NewReader(`(def x 42)`))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if _, err := Eval(expr[0], env); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if _, err := env.Get(Symbol("x")); err != nil {
		t.Errorf("variable x not set")
	}

	expr, err = parser.Parse(strings.NewReader("x"))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	result, err := Eval(expr[0], env)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != Int(42) {
		t.Errorf("unable to read the variable")
	}
}

func TestCheckers(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Bool
	}{
		{`(nil? nil)`, true},
		{`(nil? ())`, false},
		{`(nil? 0)`, false},
		{`(nil? "")`, false},
		{`(nil? (list 1 2))`, false},
		{`(bool? true)`, true},
		{`(bool? false)`, true},
		{`(bool? nil)`, false},
		{`(bool? ())`, false},
		{`(bool? 0)`, false},
		{`(bool? "")`, false},
		{`(bool? (list 1 2))`, false},
		{`(int? nil)`, false},
		{`(int? ())`, false},
		{`(int? 0)`, true},
		{`(int? 42)`, true},
		{`(int? "")`, false},
		{`(int? (list 1 2))`, false},
		{`(float? nil)`, false},
		{`(float? ())`, false},
		{`(float? 0)`, false},
		{`(float? "")`, false},
		{`(float? (list 1 2))`, false},
		{`(float? 3.1415)`, true},
		{`(str? nil)`, false},
		{`(str? ())`, false},
		{`(str? 0)`, false},
		{`(str? "")`, true},
		{`(str? "hello")`, true},
		{`(str? (list 1 2))`, false},
		{`(list? nil)`, false},
		{`(list? ())`, true},
		{`(list? 0)`, false},
		{`(list? "")`, false},
		{`(list? "hello")`, false},
		{`(list? (list 1 2))`, true},
		{`(atom? nil)`, false},
		{`(atom? ())`, false},
		{`(atom? 0)`, true},
		{`(atom? 3.1415)`, true},
		{`(atom? true)`, true},
		{`(atom? "")`, true},
		{`(atom? (list 1 2))`, false},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		env := InitBuildin()
		result, err := Eval(expr[0], env)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}
