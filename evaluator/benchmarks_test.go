package evaluator

import "testing"

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
