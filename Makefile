default: run

.PHONY: prepare
prepare:
	go get -d -v

.PHONY: build
build: prepare
	go build -o prommer

.PHONY: run
run: build
	./prommer -target-file=./test.json

.PHONY: lint
lint:
	@if gofmt -l *.go | grep .go; then \
	  echo "^- Repo contains improperly formatted go files; run gofmt -w *.go" && exit 1; \
	  else echo "All .go files formatted correctly"; fi
	go tool vet -composites=false *.go
	#go tool vet -composites=false **/*.go

.PHONY: travis-test
travis-test: lint
	go test -v
