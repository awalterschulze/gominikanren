RUN = docker run --rm --name scheme-container --user $(shell id -u):$(shell id -g) -i -t -v $(PWD):/data scheme-image

docker-build:
	docker build -t scheme-image .

run:
	$(RUN) scheme --script ./hellokanren.scm

repl: docker-build
	$(RUN) scheme microKanren.scm microKanren-wrappers.scm microKanren-test.scm microKanren-test-programs.scm

helloworld: docker-build
	$(RUN) scheme --script ./helloworld.scm

help: docker-build
	$(RUN) scheme --help