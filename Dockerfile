# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.14 AS build

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN \
  CGO_ENABLED=0 \
  GOARCH=amd64 \
  GOOS=linux \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  go build \
    -o /coda-indexer
    
# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10

COPY --from=build /coda-indexer /bin/coda-indexer

EXPOSE 8081
CMD ["/bin/coda-indexer"]
