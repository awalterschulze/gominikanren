.PHONY: travis
travis:
	make regenerate
	make test
	make vet
	make errcheck
	make diff

.PHONY: regenerate
regenerate:
	go install github.com/awalterschulze/goderive
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
