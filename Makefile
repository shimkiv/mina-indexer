.PHONY: build test docker-build docker-push

DOCKER_IMAGE ?= figment-networks/coda-indexer

build:
	go build

test:
	go test -race ./...

docker-build:
	docker build \
		-t ${DOCKER_IMAGE} \
		-f Dockerfile \
		.

docker-push:
	docker push ${DOCKER_IMAGE}