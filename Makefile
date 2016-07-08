.PHONY: run prepare build

default: run

prepare:
	go get -d -v

build:
	go build -o docker_sd

run: prepare build
	./docker_sd -target-file /Users/tomverelst/test.json
