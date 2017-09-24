(load "microKanren.scm")
(load "microKanren-wrappers.scm")
(load "microKanren-test-programs.scm")

(test-check "equalo"
  (let 
    (
      ($ 
        (
          (call/fresh 
            (lambda (q) (== q 5))
          ) 
          empty-state
        )
      )
    )
    (car $)
  )
  '(((#(0) . 5)) . 1) ;; var zero is assigned to five and the var counter is 1
)

(define (caro l a)
  (call/fresh
    (lambda (d)
      (== `(,a . ,d) l))))

(test-check "caro literal"
  (let 
    (
      ($ 
        (
          (caro `(a c o r n) `a) 
          empty-state
        )
      )
    )
    (car $)
  )
  '(((#(0) c o r n)) . 1) ;; var 0 = `(c o r n)
)

(test-check "caro fst var"
  (let 
    (
      ($ 
        (
          (call/fresh
            (lambda (q)
              (caro `(,q c o r n) `a)
            )
          )
          empty-state
        )
      )
    )
    (car $)
  )
  '(((#(1) c o r n) (#(0) . a)) . 2) ;; var 1 = `(c o r n), var 0 = 'a
)

(test-check "caro car var"
  (let 
    (
      ($ 
        (
          (call/fresh
            (lambda (q)
              (caro `(a c o r n) q)
            )
          )
          empty-state
        )
      )
    )
    (car $)
  )
  '(((#(1) c o r n) (#(0) . a)) . 2) ;; var 1 = `(c o r n), var 0 = 'a
)

(test-check "cdr 0"
  (cdr '(a))
  '()
)

(test-check "cdr 1"
  (cdr '(a b))
  '(b)
)

(test-check "cdr 2"
  (cdr '(a b c))
  '(b c)
)

(test-check "appendo single item"
  (run* (q)
    (appendo '(a) 'b q)
  )
  '((a . b))
)

(test-check "appendo single list"
  (run* (q)
    (appendo '(a) '(b) q)
  )
  '((a b))
)

(test-check "appendo combo"
  (run* (q) (fresh (x y) (== `(,x ,y) q) (appendo x y '(1 2 3 4 5))))
  '((() (1 2 3 4 5))
    ((1) (2 3 4 5))
    ((1 2) (3 4 5))
    ((1 2 3) (4 5))
    ((1 2 3 4) (5))
    ((1 2 3 4 5) ())))

;; example of unify with debugging print outs
(define (unifyWithDebugging u v s)
(begin
(display "unify:")
(display "u=")
(display u)
(display ";v=")
(display v)
(display ";s=")
(display s)
(newline)
  (let ((u (walk u s)) (v (walk v s)))
    (cond
      ((and (var? u) (var? v) (var=? u v)) s)
      ((var? u) (ext-s u v s))
      ((var? v) (ext-s v u s))
      ((and (pair? u) (pair? v))
        (begin
          (display "unify pairs:")
          (display "u=")
          (display u)
          (display ";v=")
          (display v)
          (display ";s=")
          (display s)
          (newline)
          (let ((s (unify (car u) (car v) s)))
            (and s (unify (cdr u) (cdr v) s)))
        )
      )
      (else (and (eqv? u v) s))))
)
)