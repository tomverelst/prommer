.PHONY: run prepare build

default: run

prepare:
	go get -d -v

build:
	go build -o prommer

run: prepare build
	./prommer
