# gol: Simple LISP implemented in Go

**gol** is a simple LISP with a Clojure-like syntax.  My goal was to learn more about programming language
internals, *Go, and LISP*. It all started with reading the classic [*Structure and Interpretation of Computer Programs*][sicp]
by Abelson, Sussman, and Sussman. I decided I want to write my own LISP to go even deeper into the rabbit hole.
I found the [*Build Your Own LISP*][build-lisp] book by Daniel Holden that teaches you C while building a LISP
interpreter. However I didn't like the idea of doing it in C. Nonetheless, the idea of learning another language
while building an interpreter sounded interesting, so I found the great [*Writing An Interpreter In Go*][interpreter-go]
by Thorsten Ball. To catch up with Go, I've read also [*Learning Go*][learn-go] by Jon Bodner, which was very
helpful for me. While working on it, I also found the great [Make a LISP][mal] repository that helped me to
structure my work.

While building **gol** I decided to make small changes as compared to Clojure syntax, for example, I didn't write the
`cons` function that prepends an element to a list, but instead, there is an `append` function, since internally
I'm using Go's slices for lists.

 [sicp]: https://www.goodreads.com/book/show/43713.Structure_and_Interpretation_of_Computer_Programs
 [build-lisp]: http://buildyourownlisp.com/
 [interpreter-go]: https://interpreterbook.com/
 [learn-go]: https://www.goodreads.com/book/show/55841848
 [mal]: https://github.com/kanaka/mal/
