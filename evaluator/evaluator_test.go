package evaluator

import (
	"testing"

	"github.com/google/go-cmp/cmp"

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
		{`(do (+ 2 2) (+ 3 5))`, Int(8)},
		{`'(1 2 3)`, List{Int(1), Int(2), Int(3)}},
		{`'foo`, Symbol("foo")},
		{`(append '(1 2) 3)`, List{Int(1), Int(2), Int(3)}},
		{`(append '(1) 2)`, List{Int(1), Int(2)}},
		{`(append '(1) 2 3)`, List{Int(1), Int(2), Int(3)}},
		{`(append '(1) '(2 3))`, List{Int(1), List{Int(2), Int(3)}}},
		{`(concat '(1 2) '(3 4))`, List{Int(1), Int(2), Int(3), Int(4)}},
		{`(eval (+ 2 2))`, Int(4)},
		{`(eval '(+ 2 2))`, Int(4)},
		{`(eval (list + 2 2))`, Int(4)},
		{`((fn (a) a) 123)`, Int(123)},
		{`((fn (a b) (+ a b)) 1 2)`, Int(3)},
	}

	e := NewEvaluator()

	for _, tt := range testCases {
		results, err := e.Eval(tt.input)
		result := last(results)

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
		{`(str 3.14)`, String("3.14")},
		{`(int "3.14")`, Int(3)},
		{`(float "1e-5")`, Float(1e-5)},
		{`(list "Hello World!" 42 3.14)`, List{String("Hello World!"), Int(42), Float(3.14)}},
		{`(quote 3.14)`, Float(3.14)},
		{`(quote foo)`, Symbol("foo")},
		{`(quote (foo bar))`, List{Symbol("foo"), Symbol("bar")}},
		{`(first (list 1 2 3))`, Int(1)},
		{`(rest (list 1 2 3))`, List{Int(2), Int(3)}},
		{`(nth (list 1 2 3) 1)`, Int(2)},
		{`(= 2 2)`, Bool(true)},
		{`(= 2 3)`, Bool(false)},
		{`(= 2 "2")`, Bool(false)},
		{`(= (list 1 2 3) (list 1 2 3))`, Bool(true)},
		{`(= (list 1 2 3) (list 1 "2" 3))`, Bool(false)},
		{`(< 1 2)`, Bool(true)},
		{`(< 1.0 2)`, Bool(true)},
		{`(> 1 2)`, Bool(false)},
		{`(> 1.0 2)`, Bool(false)},
	}

	e := NewEvaluator()

	for _, tt := range testCases {
		results, err := e.Eval(tt.input)
		result := last(results)

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

	e := NewEvaluator()

	for _, tt := range testCases {
		results, err := e.Eval(tt.input)
		result := last(results)

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

	e := NewEvaluator()

	for _, tt := range testCases {
		results, err := e.Eval(tt.input)
		result := last(results)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}

func TestErrorFn(t *testing.T) {
	e := NewEvaluator()
	result, err := e.Eval(`(list 1 (error "ok!") 2)`)

	if err == nil {
		t.Errorf("expected error, got result: %v", result)
	}
}

func TestDefAndDel(t *testing.T) {
	e := NewEvaluator()

	if _, err := e.Eval("(def x 42)"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	results, err := e.Eval("x")
	result := last(results)
	if err != nil {
		t.Errorf("variable x not set")
	}
	if result != Int(42) {
		t.Errorf("unable to read the variable")
	}

	results, err = e.Eval("(del x)")
	result = last(results)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != Int(42) {
		t.Errorf("expected %v, got %v", Int(42), result)
	}

	_, err = e.Eval("x")
	if err == nil {
		t.Errorf("expected not to see x, got %v", result)
	}

	_, err = e.Eval("(del y)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
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
		{`(fn? fn?)`, true},
		{`(fn? +)`, true},
		{`(fn? ())`, false},
		{`(fn? nil)`, false},
	}

	e := NewEvaluator()

	for _, tt := range testCases {
		results, err := e.Eval(tt.input)
		result := last(results)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}
