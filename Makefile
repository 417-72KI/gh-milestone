.SILENT:

list: build
	gh milestone list

help: build
	gh milestone --help

build: 
	cd cmd/gh-milestone && go build
	mv cmd/gh-milestone/gh-milestone .

install: build
	gh extension install .
