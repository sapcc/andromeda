FROM golang:1.19-alpine as builder
RUN apk add --no-cache make gcc musl-dev protoc git && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && go install go-micro.dev/v4/cmd/protoc-gen-micro@v4

COPY . /src/github.com/sapcc/andromeda
RUN make -C /src/github.com/sapcc/andromeda && mkdir -p /pkg/bin && cp /src/github.com/sapcc/andromeda/bin/* /pkg/bin/

################################################################################

FROM alpine:3.16
LABEL source_repository="https://github.com/sapcc/andromeda"

RUN apk add --no-cache ca-certificates
COPY --from=builder /pkg/ /usr/
ENTRYPOINT [ "/usr/bin/andromeda" ]
