# syntax = docker/dockerfile:1.3
FROM golang:1.16-alpine as builder

ENV GOPATH=/go
ENV GOCACHE=/gocache
ENV GOROOT=/usr/local/go
ENV CGO_ENABLED=1
ENV HELLO_WORLD_PATH=$GOROOT/src/hello-world.go

COPY ./go.mod ./go.sum ${HELLO_WORLD_PATH}/

WORKDIR ${HELLO_WORLD_PATH}

RUN go mod download

COPY ./cmd ${HELLO_WORLD_PATH}/cmd
COPY ./pkg ${HELLO_WORLD_PATH}/pkg

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/gocache \
    go mod vendor \
    && go build -v -o /bin/hello-world ./cmd/main.go

RUN if [ "$(go fmt ./... | wc -l)" -gt 0 ]; then echo "Invalid code-style. Please run 'go fmt ./...'" && exit 1; fi

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/gocache \
    go test ./...

FROM scratch

WORKDIR /

COPY --from=builder /bin/hello-world /hello-world

ENTRYPOINT ["/hello-world"]
