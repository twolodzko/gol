# Simple LISP implemented in Go

**gol** is a simple LISP with a Clojure-like syntax.  My goal was to learn more about programming language
internals, *Go, and LISP*. It all started with reading the classic [*Structure and Interpretation of Computer Programs*][sicp]
by Abelson, Sussman, and Sussman. I decided I want to write my LISP, to go even deeper into the rabbit hole.
I found the [*Build Your Own LISP*][build-lisp] book by Daniel Holden, which teaches you C while building a LISP
interpreter. However I didn't like the idea of doing it in C. Nonetheless, the idea of learning another language
while building an interpreter sounded interesting. This is how I found the great [*Writing An Interpreter In Go*][interpreter-go]
by Thorsten Ball. To catch up with Go, I've read the [*Learning Go*][learn-go] book by Jon Bodner, which was very
helpful for me. While working on it, I also found the [*Make a LISP*][mal] repository that helped me with structuring
my work. Another resource worth mentioning is the [*(How to Write a (Lisp) Interpreter (in Python))*][lispy] article
by Peter Norvig.

## Features

 * It has only the basic `int`, `float`, `string`, and `list` data types.
 * Booleans are represented as `true` and `false`.
   [As in Clojure][clj-bool], and unlike Scheme, everything except `false` and `nil` is true.
 * Values can be assigned to symbols using: `(def x 42)`.
 * Anonymous, [first-class][first-class] functions use the syntax: `(fn (x y) (+ x y))`. They can be named
   with `def`.
 * Contexts handling with `let` uses [Clojure's syntax][clj-let]: `(let (x 2 y (+ x 1)) (/ x y))`.
 * `begin`, `apply`, `map` work as in [Scheme][scheme].
 * `if` and `cond` conditionals are available, e.g. `(cond (false "not this") (true "this!"))`.
 * Lists are internally Go's [slices][go-slice], so `append` is preferred to using `cons`. Lists can be
   concatenated using `concat`. Their elements are accessed using `first`, `rest`, `init`, `last`, and `nth`. 
 * Arithmetic operators: `+`, `-`, `*`, `/`, `%` (modulo) do floating point computations and internally
   convert `int` values to `float`. If you want to do integer arithmetics, use the `int+`, `int-`, `int*`,
   `int/`, `int%` counterparts. Additionally, most of the functions from Go's [math][go-math] package
   are available under the lowercase names.
 * `quote` (`'`), `quasiquote` (``` ` ```), `unquote` (`,`), and `eval` can be used for metaprogramming.
 * Function arguments are passed by value [as in Go][pointers]. The only way to mutate a variable
   is by using `set!`.
 * Garbage collection is handled by Go's internal garbage collector.
 * [Tail call optimization][tco] is based on [github.com/kanaka/mal][mal-tco].


 [sicp]: https://www.goodreads.com/book/show/43713.Structure_and_Interpretation_of_Computer_Programs
 [build-lisp]: http://buildyourownlisp.com/
 [interpreter-go]: https://interpreterbook.com/
 [learn-go]: https://www.goodreads.com/book/show/55841848
 [clj-bool]: https://clojuredocs.org/clojure.core/boolean
 [go-math]: https://golang.org/pkg/math/
 [first-class]: https://en.wikipedia.org/wiki/First-class_function
 [go-slice]: https://blog.golang.org/slices-intro
 [clj-let]: https://clojuredocs.org/clojure.core/let
 [scheme]: https://www.cs.cmu.edu/Groups/AI/html/r4rs/r4rs_6.html
 [mal]: https://github.com/kanaka/mal/
 [lispy]: https://norvig.com/lispy.html
 [tco]: https://stackoverflow.com/questions/310974/what-is-tail-call-optimization
 [mal-tco]: https://github.com/kanaka/mal/blob/master/process/guide.md#step-5-tail-call-optimization
 [pointers]: https://krancour.medium.com/go-pointers-when-to-use-pointers-4f29256ddff3
