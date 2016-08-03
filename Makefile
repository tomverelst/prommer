IMAGE = tomverelst/prommer
VERSION = 0.0.5-alpha

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
	go tool vet -composites=false **/*.go

.PHONY: test
test: lint
	go test -v

.PHONY: travis-test
travis-test: lint
	go test -v

.PHONY: binary
binary: prepare test
	@echo "Compiling binary"
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o prommer main.go

.PHONY: docker
docker: binary
	@echo "Building Docker image"
	@docker build -t $(IMAGE):$(VERSION) .

.PHONY: run-docker
run-docker:
	@if ! docker images $(IMAGE) | awk '{print $$2}' | grep -q -F $(VERSION); then echo "$(IMAGE):$(VERSION) is not yet built. Please run 'make docker'"; false; fi
	docker run -P -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE):$(VERSION)

.PHONY: tag-latest
tag-latest:
	@if ! docker images $(IMAGE) | awk '{print $$2}' | grep -q -F $(VERSION); then echo "$(IMAGE):$(VERSION) is not yet built. Please run 'make docker'"; false; fi
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest

.PHONY: release
release: test tag-latest
	@if ! docker images $(IMAGE) | awk '{print $$2}' | grep -q -F $(VERSION); then echo "$(IMAGE):$(VERSION) is not yet built. Please run 'make docker'"; false; fi
	@echo "Pushing Docker image"
	@docker push $(IMAGE)
