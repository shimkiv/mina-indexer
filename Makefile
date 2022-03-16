.PHONY: setup build migrations test fmt queries docker docker-build docker-push

PROJECT      ?= mina-indexer
GIT_COMMIT   ?= $(shell git rev-parse HEAD)
GO_VERSION   ?= $(shell go version | awk {'print $$3'})
DOCKER_IMAGE ?= figmentnetworks/${PROJECT}
DOCKER_TAG   ?= latest

# Build the binary
build: queries migrations
	go build \
		-ldflags "\
			-X github.com/figment-networks/${PROJECT}/config.GitCommit=${GIT_COMMIT} \
			-X github.com/figment-networks/${PROJECT}/config.GoVersion=${GO_VERSION}"

# Install third-party tools
setup:
	go get -u github.com/jessevdk/go-assets-builder
	go get -u github.com/sosedoff/sqlembed

# Generate static migrations file
migrations:
	go-assets-builder store/migrations -p migrations -o store/migrations/migrations.go

# Embed SQL queries
queries:
	sqlembed -path=./store/queries -package=queries > ./store/queries/queries.go
	go fmt ./store/queries/queries.go

# Run tests
test:
	go test -race -cover ./...

# Format code
fmt:
	go fmt ./...

# Build a local docker image for testing
docker:
	docker build -t ${PROJECT} -f Dockerfile .

# Build a public docker image
docker-build:
	docker build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f Dockerfile \
		.
run: build
	./mina-indexer -config config.example.json -cmd=worker

# Tag and push docker images
docker-push: docker-build
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE}:latest
