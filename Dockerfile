# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.17 AS build

WORKDIR /go/src/github.com/figment-networks/mina-indexer

COPY . .

RUN go mod download

RUN go get -u github.com/jessevdk/go-assets-builder
RUN go get -u github.com/sosedoff/sqlembed

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10 AS release

WORKDIR /app

RUN addgroup --gid 1234 figment
RUN adduser --system --uid 1234 figment

COPY --from=build /go/src/github.com/figment-networks/mina-indexer/mina-indexer /app/mina-indexer

RUN chown -R figment:figment /app/mina-indexer

USER 1234
ENTRYPOINT ["/app/mina-indexer"]
