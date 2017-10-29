.PHONY: travis
travis:
	go get github.com/awalterschulze/goderive
	go get github.com/goccmack/gocc
	make regenerate
	make test
	make vet
	make errcheck
	make diff
	make lint

.PHONY: regenerate
regenerate:
	goderive ./...
	(cd sexpr && make regenerate)

.PHONY: vet
vet:
	go vet ./...

.PHONY: gofmt
gofmt:
	gofmt -l -s -w .

.PHONY: test
test:
	go test -v ./...

.PHONY: diff
diff:
	git diff --exit-code .

.PHONY: errcheck
errcheck:
	go get github.com/kisielk/errcheck
	errcheck ./...

lint:
	go get github.com/golang/lint/golint
	golint -set_exit_status ./micro
	golint -set_exit_status ./mini
	golint -set_exit_status ./sexpr
