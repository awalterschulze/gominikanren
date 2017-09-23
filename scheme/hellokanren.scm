(load "microKanren.scm")
(load "microKanren-wrappers.scm")
(load "microKanren-test-programs.scm")

(test-check "my own test"
  (let (($ (a-and-b empty-state)))
    (car $))
  '(((#(1) . 5) (#(0) . 7)) . 2))