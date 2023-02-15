.SILENT:

list: build
	gh milestones list

help: build
	gh milestones --help

build: 
	cd cmd/gh-milestones && go build
	mv cmd/gh-milestones/gh-milestones .