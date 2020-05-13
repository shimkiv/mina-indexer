.PHONY: build test docker docker-build docker-push

GIT_COMMIT ?= $(shell git rev-parse HEAD)
GO_VERSION ?= $(shell go version | awk {'print $$3'})
DOCKER_IMAGE ?= figment-networks/coda-indexer
DOCKER_TAG ?= latest

# Build the binary
build:
	go build \
		-ldflags "\
			-X github.com/figment-networks/coda-indexer/cli.gitCommit=${GIT_COMMIT} \
			-X github.com/figment-networks/coda-indexer/cli.goVersion=${GO_VERSION}"

# Run tests
test:
	go test -race ./...

# Build a local docker image for testing
docker:
	docker build -t coda-indexer -f Dockerfile .

# Build a public docker image
docker-build:
	docker build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f Dockerfile \
		.

# Push docker images
docker-push:
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE}:latest