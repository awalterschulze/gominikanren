.PHONY: ci
ci:
	go install github.com/goccmack/gocc
	make regenerate
	make test
	make vet
	make errcheck
	make diff
	make lint

.PHONY: regenerate
regenerate:
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
	errcheck ./...

.PHONY: lint
lint:
	golint -set_exit_status ./micro
	golint -set_exit_status ./mini
	golint -set_exit_status ./sexpr

.PHONY: bench
bench:
	go test ./... -bench=. -benchtime=10s
