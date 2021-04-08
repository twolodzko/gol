package evaluator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRecursion(t *testing.T) {
	e := NewEvaluator()

	code := `
	(def fibo (fn (n)
		(if (= n 0) 0
			(if (= n 1) 1
				(+ (fibo (- n 1))
				   (fibo (- n 2)))))))
	`

	_, err := e.EvalString(code)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	var testCases = []struct {
		input    string
		expected Int
	}{
		{`(fibo 0)`, 0},
		{`(fibo 1)`, 1},
		{`(fibo 2)`, 1},
		{`(fibo 3)`, 2},
		{`(fibo 7)`, 13},
		{`(fibo 9)`, 34},
	}

	for _, tt := range testCases {
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

func TestTailRecursion(t *testing.T) {
	code := `(def rec (fn (n)
				(if (> n 0)
					(rec (- n 1))
					nil)))`

	e := NewEvaluator()
	_, err := e.EvalString(code)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = e.EvalString(`(rec 10)`)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// this stack overflows w/o tail call optimization
	_, err = e.EvalString(`(rec 1000000)`)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func BenchmarkFibonacciRecursive(b *testing.B) {
	code := `
	(def fib (fn (n)
		(if (< n 2)
			n
			(+ (fib (- n 1))
			   (fib (- n 2))))))
	`

	e := NewEvaluator()
	e.EvalString(code)

	for i := 0; i < b.N; i++ {
		e.EvalString(`(fib 5)`)
		e.EvalString(`(fib 10)`)
		e.EvalString(`(fib 20)`)
	}
}

func BenchmarkFibonacciTailRecursive(b *testing.B) {
	code := `
	(def fib (fn (n)
		(let (loop (fn (a b n)
			(if (= n 0) a
				(loop b (+ a b) (- n 1)))))
		(loop 0 1 n))))
	`

	e := NewEvaluator()
	e.EvalString(code)

	for i := 0; i < b.N; i++ {
		e.EvalString(`(fib 5)`)
		e.EvalString(`(fib 10)`)
		e.EvalString(`(fib 20)`)
	}
}
