FROM golang:1.17-alpine as builder
RUN apk add --no-cache make gcc musl-dev protoc git && go get github.com/golang/protobuf/protoc-gen-go && go install github.com/asim/go-micro/cmd/protoc-gen-micro/v3@latest

COPY . /src/github.com/sapcc/andromeda
RUN make -C /src/github.com/sapcc/andromeda && mkdir -p /pkg/bin && cp /src/github.com/sapcc/andromeda/bin/* /pkg/bin/

################################################################################

FROM alpine:3.15
LABEL source_repository="https://github.com/sapcc/andromeda"

RUN apk add --no-cache ca-certificates
COPY --from=builder /pkg/ /usr/
ENTRYPOINT [ "/usr/bin/andromeda" ]
