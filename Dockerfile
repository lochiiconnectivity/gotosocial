# syntax=docker/dockerfile:1.3

# stage 1: Initialise ARG and defaults
ARG BUILDPLATFORM=linux/amd64
ARG TARGETPLATFORM=linux/amd64
ARG VERSION=0

# stage 2: build the go binary
FROM --platform=${BUILDPLATFORM} golang:1.19.3-alpine AS golang
COPY go.mod /go/src/github.com/superseriousbusiness/gotosocial/go.mod
COPY go.sum /go/src/github.com/superseriousbusiness/gotosocial/go.sum
COPY cmd /go/src/github.com/superseriousbusiness/gotosocial/cmd
COPY internal /go/src/github.com/superseriousbusiness/gotosocial/internal
COPY docs /go/src/github.com/superseriousbusiness/gotosocial/docs
COPY testrig /go/src/github.com/superseriousbusiness/gotosocial/testrig
WORKDIR /go/src/github.com/superseriousbusiness/gotosocial
RUN CGO_ENABLED=0 \
    go build -trimpath \
    -tags "netgo osusergo static_build kvformat" \
    -ldflags="-s -w -extldflags '-static' -X 'main.Version=${VERSION}'" \
    ./cmd/gotosocial

# stage 3: generate up-to-date swagger
FROM --platform=${BUILDPLATFORM} quay.io/goswagger/swagger:v0.30.0 AS swagger
COPY go.mod /go/src/github.com/superseriousbusiness/gotosocial/go.mod
COPY go.sum /go/src/github.com/superseriousbusiness/gotosocial/go.sum
COPY cmd /go/src/github.com/superseriousbusiness/gotosocial/cmd
COPY internal /go/src/github.com/superseriousbusiness/gotosocial/internal
WORKDIR /go/src/github.com/superseriousbusiness/gotosocial
RUN swagger generate spec -o /go/src/github.com/superseriousbusiness/gotosocial/swagger.yaml --scan-models

# stage 4: generate the web/assets/dist bundles
FROM --platform=${BUILDPLATFORM} node:16.15.1-alpine3.15 AS bundler

COPY web web
RUN yarn install --cwd web/source && \
    BUDO_BUILD=1 node web/source  && \
    rm -r web/source

# stage 5: build the executor container
FROM --platform=${TARGETPLATFORM} alpine:3.15.4 as executor

# copy from the relevant containers
COPY --chown=1000:1000 --from=golang /go/src/github.com/superseriousbusiness/gotosocial/gotosocial /gotosocial/gotosocial
COPY --chown=1000:1000 --from=swagger /go/src/github.com/superseriousbusiness/gotosocial/swagger.yaml web/assets/swagger.yaml
COPY --chown=1000:1000 --from=bundler web /gotosocial/web

WORKDIR "/gotosocial"
ENTRYPOINT [ "/gotosocial/gotosocial", "server", "start" ]
