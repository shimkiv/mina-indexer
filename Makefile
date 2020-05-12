.PHONY: build test docker-build

DOCKER_IMAGE ?= figment-networks/coda-indexer
GO_VERSION = $(shell go version | awk {'print $$3'})
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

build:
	go build

test:
	go test -race ./...

docker-build:
	docker build \
		--build-arg=GIT_COMMIT=${GIT_COMMIT} \
		--build-arg=GIT_BRANCH=${GIT_BRANCH} \
		-t ${DOCKER_IMAGE} \
		-f Dockerfile \
		.