;; binary search

(def l '(1 1 1 2 3 3 4 5 5 5 6 9 14 14 15 20 21 22 29 40 51 52 55 59 60 80 100 110 111 120 130 150 200))

(def (search val lst)
    ;; (def step 0)
    (def (iter start end)
        ;; (set! step (int+ step 1))
        (let (mid (int/ (int+ start end) 2))
            ;; (println
            ;;     step
            ;;     (list start mid end)
            ;;     (list (nth l start) (nth l mid) (nth l end)))
            (cond
                ((= start end)
                    (if (= (nth lst start) val) start nil))
                ((> val (nth lst mid))
                    (iter (int+ mid 1) end))
                ((< val (nth lst mid))
                    (iter start (int- mid 1)))
                (true mid))))
    (iter 0 (int- (count lst) 1)))

(search 60 l)

(def (assert expr)
    (if (true? expr) nil
        (error "assertion error")))

(assert (nil? (search 0 l)))
(assert (nil? (search 300 l)))
(assert (nil? (search 23 l)))
(assert (= (search 2 l) 3))
(assert (= (search 200 l) 32))
(assert (= (search 130 l) 30))
