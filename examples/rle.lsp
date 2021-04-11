;; run-length encoding

(def l '(a a b a a a a b b a c c c a a a c c b b a))

(println "Input: " l)

(def rle
    (fn (in prev acc)
	    (cond
            ((empty? in)
                (rest (conj acc prev)))
	        ((= (first in) (first prev))
                (rle
                    (rest in)
                    (list (first prev) (+ 1 (first (rest prev))))
                    acc))
	        (true
                (rle
                    (rest in)
                    (list (first in) 1)
                    (conj acc prev))))))

(println "Result:" (rle l '() '()))
