.PHONY: ci
ci:
	go install github.com/goccmack/gocc
	go install github.com/awalterschulze/goderive
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

.PHONY: lint
lint:
	go get golang.org/x/lint/golint
	golint -set_exit_status ./micro
	golint -set_exit_status ./mini
	golint -set_exit_status ./sexpr

.PHONY: bench
bench:
	go test ./... -bench=. -benchtime=10s
