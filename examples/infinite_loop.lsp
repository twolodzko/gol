
(begin
    (println "Woohoo!")
    (eval (parse-string (read-file "examples/infinite_loop.lsp"))))
