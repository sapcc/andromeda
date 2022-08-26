FROM golang:1.18-alpine as builder
RUN apk add --no-cache make gcc musl-dev protoc git && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

COPY . /src/github.com/sapcc/andromeda
RUN make -C /src/github.com/sapcc/andromeda && mkdir -p /pkg/bin && cp /src/github.com/sapcc/andromeda/bin/* /pkg/bin/

################################################################################

FROM alpine:3.16
LABEL source_repository="https://github.com/sapcc/andromeda"

RUN apk add --no-cache ca-certificates
COPY --from=builder /pkg/ /usr/
ENTRYPOINT [ "/usr/bin/andromeda-server" ]
