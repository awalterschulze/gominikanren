.PHONY: ci
ci:
	go install github.com/goccmack/gocc
	make regenerate
	make test
	make vet
	make diff

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
	go test -v -count=1 ./...

.PHONY: diff
diff:
	git diff --exit-code .

.PHONY: bench
bench:
	go test ./... -bench=. -benchtime=10s
