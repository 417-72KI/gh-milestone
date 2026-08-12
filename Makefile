.SILENT:

build: 
	cd cmd/gh-milestone && go build
	mv cmd/gh-milestone/gh-milestone .

list: build
	gh milestone list

help: build
	gh milestone --help

install: build
	gh extension install .

release: install
	gh release create $(shell gh milestone list --json title --jq '.[].title' | peco)
