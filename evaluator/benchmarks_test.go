package evaluator

import "testing"

func BenchmarkSumGol(b *testing.B) {
	e := NewEvaluator()

	for i := 0; i < b.N; i++ {
		e.EvalString(`(int+ 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20)`)
	}
}

func BenchmarkSumGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func(arr []int) int {
			tot := 0
			for _, x := range arr {
				tot += x
			}
			return tot
		}([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	}
}

func BenchmarkRecursiveGol(b *testing.B) {
	fn := `
	(def recur (fn (i n tot)
		(if (= i n) tot
			(recur (+ i 1) n (+ i tot)))))
	`

	e := NewEvaluator()
	e.EvalString(fn)

	for i := 0; i < b.N; i++ {
		e.EvalString(`(recur 0 100 0)`)
	}
}

func benchmarkRecursiveGoFn(i int, n int, tot int) int {
	if i == n {
		return tot
	}
	return benchmarkRecursiveGoFn(i+1, n, tot+1)
}

func BenchmarkRecursiveGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkRecursiveGoFn(0, 100, 0)
	}
}

func BenchmarkFactorialGol(b *testing.B) {
	fn := `
	(def fact (fn (n)
		(if (= n 1) 1
			(* n (fact (- n 1))))))
	`

	e := NewEvaluator()
	e.EvalString(fn)

	for i := 0; i < b.N; i++ {
		e.EvalString(`(fact 100.0)`)
	}
}

func benchmarkFactorialGoFn(n float64) float64 {
	if n == 1 {
		return 1
	}
	return n * benchmarkFactorialGoFn(n-1)
}

func BenchmarkFactorialGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkFactorialGoFn(100)
	}
}

func BenchmarkFibonacciGol(b *testing.B) {
	fn := `
	(def (fibo n)
		(if (= n 0) 0
			(if (= n 1) 1
				(+ (fibo (- n 1))
				   (fibo (- n 2))))))
	`
	e := NewEvaluator()
	e.EvalString(fn)

	for i := 0; i < b.N; i++ {
		e.EvalString(`(fibo 10)`)
	}
}

func benchmarkFibonacciGoFn(n float64) float64 {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return benchmarkFibonacciGoFn(n-1) + benchmarkFibonacciGoFn(n-2)
	}
}

func BenchmarkFibonacciGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkFibonacciGoFn(10)
	}
}
