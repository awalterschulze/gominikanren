package mini

import (
    "github.com/awalterschulze/gominikanren/micro"
)

/*
(defrel (very-recursiveo)
  (conde
    ((nevero))
    ((very-recursiveo))
    ((alwayso))
    ((very-recursiveo))
    ((nevero))))
*/
func VeryRecursiveO() micro.Goal {
	return Conde(
		micro.NeverO(),
		Zzz2(VeryRecursiveO),
		micro.AlwaysO(),
		Zzz2(VeryRecursiveO),
		micro.NeverO(),
	)
}
