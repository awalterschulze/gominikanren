regenerate:
	go install github.com/awalterschulze/goderive
	goderive ./...
	(cd sexpr && make regenerate)

test:
	go test -v ./...