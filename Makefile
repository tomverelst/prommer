default: run

.PHONY: prepare
prepare:
	go get -d -v

.PHONY: build
build:
	go build -o prommer

.PHONY: run
run: prepare build
	./prommer -target-file=./test.json
