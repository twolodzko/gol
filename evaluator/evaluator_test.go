package evaluator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		{`'(1 2 3)`, List{Int(1), Int(2), Int(3)}},
		{`'foo`, Symbol("foo")},
		{`(if true 1 2)`, Int(1)},
		{`(if false 1 2)`, Int(2)},
		{`(if (true? false)
			  (error "this should not fail!")
			  "ok")`, String("ok")},
		{`(cond ((= 2 1) "wrong"))`, nil},
		{`(cond
			((= 2 1) "wrong")
			(true "correct"))`, String("correct")},
		{`(cond
			(false 1)
			(true 2)
			(true (error "Oh, no!")))`, Int(2)},
		{`(quote (+ 1 2))`, List{Symbol("+"), Int(1), Int(2)}},
		{`(quasiquote (unquote (+ 1 2)))`, Int(3)},
		{"`,(+ 1 2)", Int(3)},
		{`(def x 4)
		  (quasiquote
			(+ 1 2 (unquote (+ 1 2)) (unquote x)))`, List{Symbol("+"), Int(1), Int(2), Int(3), Int(4)}},
		{"(def x 4) `(+ 1 2 ,(+ 1 2) (- ,x))", List{Symbol("+"), Int(1), Int(2), Int(3), List{Symbol("-"), Int(4)}}},
		{"``,x", List{Symbol("quasiquote"), List{Symbol("unquote"), Symbol("x")}}},
		{"(def x 5) (eval (eval ```,x))", Int(5)},
		{`(eval '(+ 2 2))`, Int(4)},
		{`(- 7 (* 2 (+ 1 2)) 1)`, Int(0)},
		{`(def b (+ 1 2)) b`, Int(3)},
		{`(let (c 2) c)`, Int(2)},
		{`(let (x 10)
			(+ 5 x))`, Int(15)},
		{`(let (x 5)
			(let (y 6)
				 (+ x y)))`, Int(11)},
		{`(let (x 1 y (+ 2 3))
			(+ x y))`, Int(6)},
		{`(let (x 1 y (+ 1 x))
			(+ x y))`, Int(3)},
		{`(begin (+ 2 2)
			     (+ 3 5))`, Int(8)},
		{`(begin (def x 2)
			     (+ x 4))`, Int(6)},
		{`(eval (+ 2 2))`, Int(4)},
		{`(eval '(+ 2 2))`, Int(4)},
		{`(eval (list + 2 2))`, Int(4)},
		{`((fn (a) a) 123)`, Int(123)},
		{`((fn (a b)
			(+ a b))
		  1 2)`, Int(3)},
		{`((fn (x y)
			(or (= x y)
				(> x y)))
		  2 1)`, Bool(true)},
		{`((fn () 5))`, Int(5)},
		{`(def x 2)
		  ((fn (x)
		  	(+ x 5))
			(+ x 3))`, Int(10)},
		{`(def factorial (fn (n)
			(if (= n 1) 1
				(* n (factorial (- n 1))))))
		  (factorial 10)`, Int(3628800)},
		{`(def foo (fn (x)
			(fn (y) (+ x y))))
		  ((foo 5) 6)`, Int(11)},
		{`(def foo (fn (x)
			(+ x 1)))
		  (let (x foo)
		  	(let (y x)
			  (y 4)))`, Int(5)},
		{`(def x nil)
		  (let () (set! x 4))
		  x`, Int(4)},
	}

	for _, tt := range testCases {
		e := NewEvaluator()
		results, err := e.EvalString(tt.input)
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
		{`(init (list 1 2 3))`, List{Int(1), Int(2)}},
		{`(last (list 1 2 3))`, Int(3)},
		{`(nth (list 1 2 3) 1)`, Int(2)},
		{`(append '(1 2) 3)`, List{Int(1), Int(2), Int(3)}},
		{`(append '(1) 2)`, List{Int(1), Int(2)}},
		{`(append '(1) 2 3)`, List{Int(1), Int(2), Int(3)}},
		{`(append '(1) '(2 3))`, List{Int(1), List{Int(2), Int(3)}}},
		{`(cons 1 '(2 3 4))`, List{Int(1), Int(2), Int(3), Int(4)}},
		{`(cons '(1 2) '(3 4))`, List{List{Int(1), Int(2)}, Int(3), Int(4)}},
		{`(concat '(1 2) '(3 4))`, List{Int(1), Int(2), Int(3), Int(4)}},
		{`(= 2 2)`, Bool(true)},
		{`(= 2 2 2)`, Bool(true)},
		{`(= 2 3)`, Bool(false)},
		{`(= 2 2 3)`, Bool(false)},
		{`(= 2 3 2)`, Bool(false)},
		{`(= 2 "2")`, Bool(false)},
		{`(= (list 1 2 3) (list 1 2 3))`, Bool(true)},
		{`(= (list 1 2 3) (list 1 "2" 3))`, Bool(false)},
		{`(< 1 2)`, Bool(true)},
		{`(< 1.0 2)`, Bool(true)},
		{`(> 1 2)`, Bool(false)},
		{`(> 1.0 2)`, Bool(false)},
		{`(< 1 2 3)`, Bool(true)},
		{`(< 1 2 2 3)`, Bool(false)},
		{`(> 3 2 1)`, Bool(true)},
		{`(> 3 2 1 1)`, Bool(false)},
		{`(eval (parse-string "(+ 2 2)"))`, Int(4)},
		{`(parse-string (str '(1 2 "3")))`, List{Int(1), Int(2), String("3")}},
		{`(apply (fn (x) x) '('test))`, Symbol("test")},
		{`(apply + '(1 2 3))`, Int(6)},
		{`(map (fn (x) x) '(1 2 3))`, List{Int(1), Int(2), Int(3)}},
		{`(map - '(1 2 3))`, List{Int(-1), Int(-2), Int(-3)}},
	}

	for _, tt := range testCases {
		e := NewEvaluator()
		results, err := e.EvalString(tt.input)
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

	for _, tt := range testCases {
		e := NewEvaluator()
		results, err := e.EvalString(tt.input)
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

	for _, tt := range testCases {
		e := NewEvaluator()
		results, err := e.EvalString(tt.input)
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
	result, err := e.EvalString(`(list 1 (error "ok!") 2)`)

	if err == nil {
		t.Errorf("expected error, got result: %v", result)
	}
}

func TestDef(t *testing.T) {
	e := NewEvaluator()

	if _, err := e.EvalString("(def x 42)"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	results, err := e.EvalString("x")
	result := last(results)
	if err != nil {
		t.Errorf("variable x not set")
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
		{`(empty? '())`, true},
		{`(empty? '(1 2))`, false},
		{`(empty? (list))`, true},
	}

	for _, tt := range testCases {
		e := NewEvaluator()
		results, err := e.EvalString(tt.input)
		result := last(results)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}
